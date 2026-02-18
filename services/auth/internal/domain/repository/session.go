package repository

import (
	"auth/internal/domain/entity"
	"context"
)

type SessionRepository interface {
	// Persistence
	CreateSession(ctx context.Context, session *entity.Session) error
	CreateRefreshToken(ctx context.Context, session *entity.RefreshToken) error

	// Retrieval
	GetSessionByID(ctx context.Context, id string) (*entity.Session, error)
	GetSessionByAccessToken(ctx context.Context, token entity.HashedToken) (*entity.Session, error)

	GetRefreshTokenEntityByRefreshToken(ctx context.Context, token entity.HashedToken) (*entity.RefreshToken, error)

	// Updates & Security
	UpdateSession(ctx context.Context, session *entity.Session) error
	RevokeRefreshToken(ctx context.Context, sessionID string) error // For rotation detection

}
