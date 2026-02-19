package entity

import (
	"auth/internal/domain"
	"encoding/json"
	"strings"
	"time"
)

type HashedToken string
type RefreshToken struct {
	id          string
	sessionID   string
	hashedToken HashedToken
	createdAt   time.Time
	expiresAt   time.Time
	revokedAt   time.Time
}

func NewRefreshToken(id, sessionID string, hashedToken HashedToken, createdAt, expiresAt time.Time, revokedAt time.Time) (*RefreshToken, error) {

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
		expiresAt:   expiresAt,
		revokedAt:   revokedAt,
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
}

func (rt *RefreshToken) IsRevoked() bool {
	return !rt.revokedAt.IsZero()
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

type refreshTokenJSON struct {
	ID          string `json:"id"`
	SessionID   string `json:"sid"`
	HashedToken string `json:"ht"`
	CreatedAt   int64  `json:"cat"`
	ExpiresAt   int64  `json:"eat"`
	RevokedAt   int64  `json:"rat"`
}

func (rt *RefreshToken) Serialize() (string, error) {
	j := refreshTokenJSON{
		ID:          rt.id,
		SessionID:   rt.sessionID,
		HashedToken: string(rt.hashedToken),
		CreatedAt:   rt.createdAt.Unix(),
		ExpiresAt:   rt.expiresAt.Unix(),
		RevokedAt:   -1,
	}
	if rt.IsRevoked() {
		j.RevokedAt = rt.revokedAt.Unix()
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func UnserializeRefreshToken(s string) (*RefreshToken, error) {
	var j refreshTokenJSON
	if err := json.Unmarshal([]byte(s), &j); err != nil {
		return nil, err
	}
	revokedAt := time.Time{}
	if j.RevokedAt != -1 {
		revokedAt = time.Unix(j.RevokedAt, 0)
	}
	return &RefreshToken{
		id:          j.ID,
		sessionID:   j.SessionID,
		hashedToken: HashedToken(j.HashedToken),
		createdAt:   time.Unix(j.CreatedAt, 0),
		expiresAt:   time.Unix(j.ExpiresAt, 0),
		revokedAt:   revokedAt,
	}, nil
}
