package entity

import (
	"auth/internal/domain"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type HashedToken string
type RefreshToken struct {
	id          string
	sessionID   string
	hashedToken HashedToken
	isRevoked   bool
	createdAt   time.Time
	expiresAt   time.Time
	revokedAt   time.Time
}

func NewRefreshToken(id, sessionID string, hashedToken HashedToken, createdAt time.Time, expDuration time.Duration) (*RefreshToken, error) {

	id = strings.TrimSpace(id)
	sessionID = strings.TrimSpace(sessionID)
	if id == "" {
		return nil, domain.ErrEmptyId
	}

	if sessionID == "" {
		return nil, domain.ErrEmptySessionID
	}

	if hashedToken == "" {
		return nil, domain.ErrEmptyToken
	}

	return &RefreshToken{
		id:          id,
		hashedToken: hashedToken,
		sessionID:   sessionID,
		createdAt:   createdAt,
		expiresAt:   createdAt.Add(expDuration),
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

	rt.Revoke(now)

	return &RefreshToken{
		id:          id,
		hashedToken: hashedToken,
		sessionID:   rt.sessionID,
		createdAt:   now,
		expiresAt:   expiresAt,
	}, nil
}

func (rt *RefreshToken) EnsureValid(now time.Time) error {

	if rt.IsRevoked() {
		return domain.ErrReUsedToken
	}

	if now.After(rt.expiresAt) {
		rt.Revoke(now)
		return domain.ErrTokenExpired
	}

	return nil
}

func (rt *RefreshToken) Revoke(now time.Time) {
	rt.revokedAt = now
	rt.isRevoked = true
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.isRevoked
}
func (rt *RefreshToken) ID() string {
	return rt.id
}
func (rt *RefreshToken) SessionID() string {
	return rt.sessionID
}
func (rt *RefreshToken) HashedToken() HashedToken {
	return rt.hashedToken
}
func (rt *RefreshToken) CreatedAt() time.Time {
	return rt.createdAt
}
func (rt *RefreshToken) RevokedAt() time.Time {
	return rt.revokedAt
}
func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}

type SerializeRefreshTokenE string

// Format: id|sessionID|hashedToken|isRevoked|createdAt(Unix)|expiresAt(Unix)
func (rt *RefreshToken) Serialize() SerializeRefreshTokenE {
	isRevoked := "0"
	if rt.IsRevoked() {
		isRevoked = "1"
	}

	return SerializeRefreshTokenE(fmt.Sprintf("%s|%s|%s|%s|%d|%d|%d",
		rt.ID(),
		rt.SessionID(),
		rt.HashedToken(),
		isRevoked,
		rt.CreatedAt().Unix(),
		rt.ExpiresAt().Unix(),
		rt.RevokedAt().Unix(),
	),
	)
}
func Unserialize(s SerializeRefreshTokenE) *RefreshToken {
	val := string(s)
	parts := strings.Split(val, "|")
	if len(parts) < 7 {
		return nil
	}

	// Parse timestamps
	createdUnix, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return nil
	}
	expiresUnix, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return nil
	}
	revokedUnix, err := strconv.ParseInt(parts[6], 10, 64)

	// Parse boolean
	isRevoked := parts[3] == "1"

	return &RefreshToken{
		id:          parts[0],
		sessionID:   parts[1],
		hashedToken: HashedToken(parts[2]),
		isRevoked:   isRevoked,
		createdAt:   time.Unix(createdUnix, 0),
		expiresAt:   time.Unix(expiresUnix, 0),
		revokedAt:   time.Unix(revokedUnix, 0),
	}
}
