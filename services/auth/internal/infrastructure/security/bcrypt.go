package security

import (
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"auth/internal/domain/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewHasher(cost int) repository.PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (entity.HashedPassword, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", domain.ErrPasswordTooLong
		}
		return "", fmt.Errorf("generating bcrypt hash: %w", err)
	}
	return entity.HashedPassword(bytes), nil
}

func (h *BcryptHasher) CheckPasswordHash(password string, hash entity.HashedPassword) bool {
	if hash == "" || len(password) > 72 {
		// Perform dummy bcrypt to maintain constant timing
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$14$dummy"), []byte(password))
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
