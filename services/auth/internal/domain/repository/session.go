package repository

import (
	"auth/internal/domain/entity"
	"context"
)

type SessionRepository interface {
	// Persistence
	CreateSession(ctx context.Context, session *entity.Session) error
	CreateRefreshToken(ctx context.Context, session *entity.RefreshToken) error
}
