package domain

import "errors"

type ValidationError struct {
	Field   string
	Message string
}

func (v ValidationError) Error() string {
	return v.Message
}

var (
	ErrEmptyEmail    = &ValidationError{Field: "email", Message: "email is required"}
	ErrEmptyPassword = &ValidationError{Field: "password", Message: "password is required"}
	ErrEmptyName     = &ValidationError{Field: "name", Message: "name is required"}
	ErrEmptyId       = &ValidationError{Field: "id", Message: "id is required"}
)

// Auth errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongPassword      = errors.New("wrong password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Password errors

var (
	ErrTooLongPassword  = &ValidationError{Field: "password", Message: "password too long"}
	ErrPasswordTooShort = &ValidationError{Field: "password", Message: "password too short"}
	ErrPasswordTooWeak  = &ValidationError{Field: "password", Message: "password too weak"}
	ErrPasswordNotMatch = &ValidationError{Field: "password", Message: "password not match"}
)
