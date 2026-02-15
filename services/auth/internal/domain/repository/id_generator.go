package repository

type IdGenerator interface {
	NewId() string
}
