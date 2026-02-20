package helpers

import "github.com/aligh5331/godrop/services/auth/internal/domain/entity"

type PasswordHasher interface {
	Hash(password string) (entity.HashedPassword, error)
	CheckPasswordHash(password string, hashed entity.HashedPassword) bool
}

type TokenHasher interface {
	Hash(token string) (entity.HashedToken, error)
	CheckHash(token string, hashed entity.HashedToken) bool
}
