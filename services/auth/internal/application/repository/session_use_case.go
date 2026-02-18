package repository

import (
	"auth/internal/application/dto"
	"context"
)

type SessionUseCase interface {
	CreateNewSession(ctx context.Context, userID string, metadataDTO dto.SessionMetadataDTO) (*dto.TokenPairDTO, error)
	RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairDTO, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	EnsureAccessTokenValid(ctx context.Context, AccessT string) error
	GetUserSessions(ctx context.Context, userID string) ([]*dto.SessionDTO, error)
}
