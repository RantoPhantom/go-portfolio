package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"errors"
	"io/fs"
	"learning/go-portfolio/custom_errors"
	"os"
	"regexp"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

const (
	DB_PATH             = "./dbs/"
	SESSIONS_DB_PATH    = DB_PATH + "sessions.sqlite"
	TOKEN_EXPIRATION_HR = 2
)

type DB_Connection struct {
	Db      *sql.DB
	Queries *Queries
}

var (
	db_connections        = make(map[string]*DB_Connection)
	db_mutex              sync.Mutex
	session_db_connection *DB_Connection
)

// load the schemas into variables

//go:embed schemas/user_db.sql
var user_db_ddl string

//go:embed schemas/sessions_db.sql
var session_db_ddl string

func CreateSessionDB() error {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", SESSIONS_DB_PATH)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, session_db_ddl); err != nil {
		return err
	}

	queries := New(db)
	session_db_connection = &DB_Connection{
		Db:      db,
		Queries: queries,
	}
	return nil
}

// checks if the session is expired in the db and retrive the session if not
// Errors:
// - SessionExpired - the session has been expired
func CloseSessionDb() {
	session_db_connection.Db.Close()
}

func GetSessionFromToken(ctx context.Context, session_token string) (*Session, error) {
	if session_db_connection == nil {
		return nil, custom_errors.SessionDbNotInitialized
	}
	queries := session_db_connection.Queries

	session, err := queries.Get_session(ctx, session_token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		err = queries.Remove_session(ctx, session_token)
		if err != nil {
			return nil, err
		}
		return nil, custom_errors.SessionExpired
	}
	return &session, nil
}

var banned_regex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func is_valid_username(username string) bool {
	return banned_regex.Match([]byte(username))
}

// returns the user_db object from uuid
// Errors:
// UserNotFound - uuid not found in fs
// sql db errors
func GetDB(username string) (*DB_Connection, error) {
	db_mutex.Lock()
	defer db_mutex.Unlock()

	// return pre exisiting connection if exists
	if connection, ok := db_connections[username]; ok {
		return connection, nil
	}

	// check if the user is in db dir
	var user_db_path string = DB_PATH + username + ".sqlite"
	if _, err := os.Stat(user_db_path); errors.Is(err, fs.ErrNotExist) {
		return nil, custom_errors.UserNotFound
	}

	db, err := sql.Open("sqlite3", user_db_path)
	if err != nil {
		return nil, err
	}

	queries := New(db)
	user_db := &DB_Connection{
		Db:      db,
		Queries: queries,
	}

	db_connections[username] = user_db

	return db_connections[username], nil
}

// makes a random 256 bit opaque token
func generate_session_token() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func Remove_Session(ctx context.Context, session_token string) error {
	err := session_db_connection.Queries.Remove_session(ctx, session_token)
	if err != nil {
		return err
	}
	return nil
}

func CreateSession(ctx context.Context, username string) (*Session, error) {
	token, err := generate_session_token()
	if err != nil {
		return nil, err
	}
	session, err := session_db_connection.Queries.Insert_session(ctx, Insert_sessionParams{
		Username:    username,
		Token:       token,
		ExpiresAt:   time.Now().Add(TOKEN_EXPIRATION_HR * time.Hour),
		DateCreated: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func CreateDB(ctx context.Context, username string, password string) error {
	// check if main db dir exists, if not create it
	if _, err := os.Stat(DB_PATH); errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(DB_PATH, os.ModePerm)
	}

	if is_valid_username(username) {
		return custom_errors.InvalidUsername
	}
	var user_db_path string = DB_PATH + username + ".sqlite"

	db, err := sql.Open("sqlite3", user_db_path)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, user_db_ddl); err != nil {
		return err
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}

	queries := New(db)
	err = queries.Insert_user_info(ctx, Insert_user_infoParams{
		PasswordHash: string(password_hash),
		DateCreated:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExists(username string) error {
	if _, err := os.Stat(DB_PATH + username + ".sqlite"); err == nil {
		return custom_errors.UserDbExists
	}
	return nil
}
