package repository

type Validator interface {
	ValidatePassword(password string) error
	ValidateEmail(email string) error
}
