package entity

import (
	"auth/internal/domain"
	"strings"
	"time"
)

type Session struct {
	id           string
	familyID     string
	token        HashedToken
	userID       string
	userAgent    string
	ip           string
	isRevoked    bool
	createdAt    time.Time
	expiresAt    time.Time
	lastActiveAt time.Time
}

func NewSession(
	id,
	familyID,
	userAgent,
	ip,
	userID string,
	hashedToken HashedToken,
	now time.Time,
	expDuration time.Duration,
) (*Session, error) {
	id = strings.TrimSpace(id)
	familyID = strings.TrimSpace(familyID)
	userID = strings.TrimSpace(userID)

	if id == "" {
		return nil, domain.ErrEmptyId
	}
	if familyID == "" {
		return nil, domain.ErrEmptyFamilyID
	}
	if userID == "" {
		return nil, domain.ErrEmptyUserId
	}
	if hashedToken == "" {
		return nil, domain.ErrEmptyToken
	}

	return &Session{
		id:        id,
		familyID:  familyID,
		token:     hashedToken,
		userAgent: userAgent,
		ip:        ip,
		userID:    userID,
		createdAt: now.UTC(),
		expiresAt: now.Add(expDuration).UTC(),
	}, nil
}

func (s *Session) IsValid() bool {
	return !s.isRevoked && time.Now().Before(s.expiresAt)
}

func (s *Session) Revoke() {
	s.isRevoked = true
}

func (s *Session) Use(now time.Time) {
	s.lastActiveAt = now.UTC()
}
