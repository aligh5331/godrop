package security

import (
	"auth/internal/domain/entity"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SHA256TokenHasher struct{}

func (h SHA256TokenHasher) Hash(token string) (entity.HashedToken, error) {
	hash := sha256.Sum256([]byte(token))
	return entity.HashedToken(hex.EncodeToString(hash[:])), nil
}

func (h SHA256TokenHasher) CheckHash(token string, hashed entity.HashedToken) bool {
	hash := sha256.Sum256([]byte(token))
	expected := hex.EncodeToString(hash[:])
	return hmac.Equal([]byte(expected), []byte(hashed))
}
