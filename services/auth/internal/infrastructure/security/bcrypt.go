package security

import (
	"errors"
	"fmt"

	"github.com/aligh5331/godrop/services/auth/internal/domain"
	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"
	"github.com/aligh5331/godrop/services/auth/internal/domain/helpers"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewHasher(cost int) helpers.PasswordHasher {
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
