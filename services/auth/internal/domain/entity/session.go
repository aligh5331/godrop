package entity

import (
	"auth/internal/domain"
	"strings"
	"time"
)

type Session struct {
	id        string
	token     HashedToken
	userID    string
	userAgent string
	ip        string
	isRevoked bool
	createdAt time.Time
	expiresAt time.Time
}

func NewSession(
	id,
	userAgent,
	ip,
	userID string,
	hashedToken HashedToken,
	now time.Time,
	expDuration time.Duration,
) (*Session, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)

	if id == "" {
		return nil, domain.ErrEmptyId
	}
	if userID == "" {
		return nil, domain.ErrEmptyUserId
	}
	if hashedToken == "" {
		return nil, domain.ErrEmptyToken
	}

	return &Session{
		id:        id,
		token:     hashedToken,
		userAgent: userAgent,
		ip:        ip,
		userID:    userID,
		createdAt: now.UTC(),
		expiresAt: now.Add(expDuration).UTC(),
	}, nil
}

func (s *Session) IsValid(now time.Time) bool {
	return !s.isRevoked && now.Before(s.expiresAt)
}

func (s *Session) Revoke() {
	s.isRevoked = true
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Token() HashedToken {
	return s.token
}
func (s *Session) SetToken(token HashedToken, now time.Time, duration time.Duration) {
	s.token = token
	s.expiresAt = now.Add(duration)
}

func (s *Session) UserID() string {
	return s.userID
}

func (s *Session) IP() string {
	return s.ip
}
func (s *Session) UserAgent() string {
	return s.userAgent
}
