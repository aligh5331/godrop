package repository

import "auth/internal/domain/entity"

type PasswordHasher interface {
	Hash(password string) (entity.HashedPassword, error)
	CheckPasswordHash(password string, hashed entity.HashedPassword) bool
}
