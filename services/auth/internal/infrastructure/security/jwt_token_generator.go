package security

import (
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/application/dto"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secret []byte
}

func NewJWTTokenGenerator(secret []byte) JWTTokenGenerator {
	return JWTTokenGenerator{secret: secret}
}

type accessClaims struct {
	UserID      string `json:"uid"`
	IP          string `json:"ip"`
	ClientAgent string `json:"ca"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

func (j JWTTokenGenerator) GenerateAccessToken(userID string, deviceInfo dto.SessionMetadataDTO, expiration time.Duration) (string, error) {
	claims := accessClaims{
		UserID:      userID,
		IP:          deviceInfo.IP,
		ClientAgent: deviceInfo.ClientAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j JWTTokenGenerator) GenerateRefreshToken(userID string, sessionID string, expiration time.Duration) (string, error) {
	claims := refreshClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
