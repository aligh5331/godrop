package entity

import "time"

type Session struct {
	id           string
	token        string
	userID       string
	userAgent    string
	ip           string
	isRevoked    bool
	createdAt    time.Time
	expiresAt    time.Time
	lastActiveAt time.Time
}

func (s *Session) NewSession(
	id,
	token,
	userAgent,
	ip string,
	userID *User,
	now time.Time,
	expDuration time.Duration,
) {

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
