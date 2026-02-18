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
	ErrEmailNotVerified   = errors.New("email is not verified")
)

// Password errors

var (
	ErrPasswordTooLong  = &ValidationError{Field: "password", Message: "password too long"}
	ErrPasswordTooShort = &ValidationError{Field: "password", Message: "password too short"}
	ErrPasswordTooWeak  = &ValidationError{Field: "password", Message: "password too weak"}
	ErrPasswordNotMatch = &ValidationError{Field: "password", Message: "password not match"}
)

var (
	ErrReUsedToken    = errors.New("refresh token is already used")
	ErrEmptyToken     = errors.New("token is empty")
	ErrEmptySessionID = errors.New("family id is empty")
	ErrTokenExpired   = errors.New("token is expired")
	ErrRevokedToken   = errors.New("token is revoked")
	ErrEmptyUserId    = errors.New("user id is empty")
)
