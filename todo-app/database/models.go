// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"
)

type TodoItem struct {
	ID          int64
	ItemNumber  int64
	Content     string
	DateCreated time.Time
	IsDone      bool
}

type UserInfo struct {
	PasswordHash string
	DateCreated  time.Time
}
