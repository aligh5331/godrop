package entity

import (
	"auth/internal/domain"
	"time"
)

type HashedToken string
type RefreshToken struct {
	id          string
	familyID    string
	hashedToken HashedToken
	isRevoked   bool
	createdAt   time.Time
	expiresAt   time.Time
}

func NewRefreshToken(id, familyID string, hashedToken HashedToken, createdAt time.Time, expDuration time.Duration) (*RefreshToken, error) {

	if id == "" {
		return nil, domain.ErrEmptyId
	}

	if familyID == "" {
		return nil, domain.ErrEmptyFamilyID
	}

	if hashedToken == "" {
		return nil, domain.ErrEmptyToken
	}

	return &RefreshToken{
		id:          id,
		hashedToken: hashedToken,
		familyID:    familyID,
		createdAt:   createdAt.UTC(),
		expiresAt:   createdAt.Add(expDuration).UTC(),
	}, nil
}

func (rt *RefreshToken) Rotate(id string, hashedToken HashedToken, now, expiresAt time.Time) (*RefreshToken, error) {

	if err := rt.EnsureValid(now); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, domain.ErrEmptyId
	}

	if hashedToken == "" {
		return nil, domain.ErrEmptyToken
	}

	rt.Revoke()

	return &RefreshToken{
		id:          id,
		hashedToken: hashedToken,
		familyID:    rt.familyID,
		createdAt:   now.UTC(),
		expiresAt:   expiresAt.UTC(),
	}, nil
}

func (rt *RefreshToken) EnsureValid(now time.Time) error {

	if rt.IsRevoked() {
		return domain.ErrReUsedToken
	}

	if now.After(rt.expiresAt) {
		rt.Revoke()
		return domain.ErrTokenExpired
	}

	return nil
}

func (rt *RefreshToken) Revoke() {
	rt.isRevoked = true
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.isRevoked
}
func (rt *RefreshToken) ID() string {
	return rt.id
}
func (rt *RefreshToken) FamilyID() string {
	return rt.familyID
}
func (rt *RefreshToken) HashedToken() HashedToken {
	return rt.hashedToken
}
func (rt *RefreshToken) CreatedAt() time.Time {
	return rt.createdAt
}
func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}
