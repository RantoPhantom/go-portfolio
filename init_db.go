package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


func init_db() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db")
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()

	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
