package database

import (
	"time"
	"regexp"
	"sync"
	"os"
	"errors"
	"context"
	"database/sql"
	_ "embed"
	"io/fs"

	"learning/go-portfolio/custom_errors"
	_ "github.com/mattn/go-sqlite3"
)

const DB_PATH = "./dbs/"

type UserDb struct {
	Db *sql.DB
	Queries *Queries
}

var (
	db_connections = make(map[string]*UserDb)
	db_mutex sync.Mutex
)

//load the schemas into this variable
//go:embed schemas/*.sql
var ddl string

var banned_regex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func is_username_banned(username string) bool{
	return banned_regex.Match([]byte(username))
}

func GetDB(username string) (*UserDb, error) {
	db_mutex.Lock()
	defer db_mutex.Unlock()

	// return pre exisiting connection if exists
	if connection, ok := db_connections[username]; ok {
		return connection, nil
	}

	// check if the user is in db dir
	var user_db_path string = DB_PATH + username + ".sqlite"
	if _,err := os.Stat(user_db_path); errors.Is(err, fs.ErrNotExist){
		return nil, custom_errors.UserNotFound
	}

	db,err := sql.Open("sqlite3",user_db_path)
	if err != nil { return nil,err }

	queries := New(db)
	user_db := &UserDb{
		Db:db,
		Queries:queries,
	}

	db_connections[username] = user_db

	return db_connections[username], nil
}

func CreateDB(username string, password string) error {
	// check if main db dir exists, if not create it
	if _,err := os.Stat(DB_PATH); errors.Is(err, fs.ErrNotExist){
		os.Mkdir(DB_PATH, os.ModePerm)
	}

	if is_username_banned(username) {
		return custom_errors.InvalidUsername
	}
	var user_db_path string = DB_PATH + username + ".sqlite"

	ctx := context.Background()
	db,err := sql.Open("sqlite3",user_db_path)
	if err != nil { return err }

	if _, err := db.ExecContext(ctx, ddl); err != nil { return err }
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode = WAL;"); err != nil { return err }

	queries := New(db)
	err = queries.Insert_user_info(ctx, Insert_user_infoParams{
		PasswordHash: password,
		DateCreated: time.Now(),
	})
	if err != nil { return err }

	return nil
}
