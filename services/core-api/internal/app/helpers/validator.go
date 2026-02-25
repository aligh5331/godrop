package helpers

type Validator interface {
	Email(email string) error
}
