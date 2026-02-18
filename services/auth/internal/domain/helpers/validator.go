package helpers

type Validator interface {
	ValidatePassword(password string) error
	ValidateEmail(email string) error
	ValidateIP(ip string) error
}
