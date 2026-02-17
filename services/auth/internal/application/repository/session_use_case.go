package repository

import (
	"auth/internal/application/dto"
	"context"
)

type SessionUseCase interface {
	CreateNewSession(ctx context.Context, userID string) (*dto.TokenPairDTO, error)
	RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairDTO, error)
	RevokeSession(ctx context.Context, userID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	GetUserSessions(ctx context.Context, userID string) error
}
