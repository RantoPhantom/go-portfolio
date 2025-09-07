package custom_errors

import (
	"errors"
)

var (
	UserNotFound            = errors.New("This user does not exist.")
	InvalidInput            = errors.New("This input cannot be empty.")
	UserDbExists            = errors.New("This user is already registered.")
	InvalidUsername         = errors.New("This username is not usable")
	InvalidCredentials      = errors.New("Invalid username or password")
	PasswordTooLong         = errors.New("This password exceeds the limit!")
	SessionDbExists         = errors.New("Session Db alr exists")
	SessionExpired          = errors.New("This session has expired, please log in again.")
	SessionDbNotInitialized = errors.New("Session db has not been initialized.")
	EnvEmpty                = errors.New("env var: %s has not been initialized.")
)
