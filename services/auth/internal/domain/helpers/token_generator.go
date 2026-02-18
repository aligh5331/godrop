package helpers

import (
	"auth/internal/application/dto"
	"time"
)

type TokenGenerator interface {
	GenerateAccessToken(userID string, deviceInfo dto.SessionMetadataDTO, expiration time.Duration) (string, error)
	GenerateRefreshToken(userID string, sessionID string, expiration time.Duration) (string, error)
}
