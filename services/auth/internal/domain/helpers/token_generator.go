package helpers

import (
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/application/dto"
)

type TokenGenerator interface {
	GenerateAccessToken(userID string, deviceInfo dto.SessionMetadataDTO, expiration time.Duration) (string, error)
	GenerateRefreshToken(userID string, sessionID string, expiration time.Duration) (string, error)
}
