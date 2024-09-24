package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


func init_db() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	_,createErr := db.Exec("CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY AUTOINCREMENT, content TEXT)")
	if createErr != nil {
		return nil, createErr
	}

	return db, nil
}
