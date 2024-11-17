package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init_db()  {
	var err error

	db, err = sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = init_tables()
	if err != nil{
		log.Fatal(err)
	}
}

func GetDB() (*sql.DB) {
	if db == nil{
		init_db()
	}
	return db
}

func init_tables() error {
	var create_query string = ""
	var err error = nil

	//init the todo table
	create_query = `CREATE TABLE IF NOT EXISTS todo (
	id INTEGER NOT NULL,
	todo_content TEXT,
	date_created TEXT,
	is_done INTEGER DEFAULT FALSE)
	`
	_,err = db.Exec(create_query)
	if err != nil {
		return err
	}

	//init the users table
	create_query = `CREATE TABLE IF NOT EXISTS chat_users (
	user_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	date_created TEXT)
	`
	_,err = db.Exec(create_query)
	if err != nil {
		return err
	}

	return nil
}
