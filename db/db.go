package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init_db() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return db, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	create_table := "CREATE TABLE IF NOT EXISTS todo (id INTEGER NOT NULL, todo_content TEXT, date_created TEXT, is_done INTEGER DEFAULT FALSE)"
	_,createErr := db.Exec(create_table)
	if createErr != nil {
		return nil, createErr
	}

	return db, nil
}

func GetDB() (*sql.DB, error) {
	var err error
	if db == nil{
		db, err = Init_db()
		if err != nil { return nil, err}
	}
	return db, nil
}
