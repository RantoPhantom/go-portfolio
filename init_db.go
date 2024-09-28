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

	create_table := "CREATE TABLE IF NOT EXISTS todo (id INTEGER NOT NULL, todo_content TEXT, date_created TEXT, is_done INTEGER DEFAULT FALSE)"
	_,createErr := db.Exec(create_table)
	if createErr != nil {
		return nil, createErr
	}

	return db, nil
}
