package repository

import (
	"context"

	"github.com/aligh5331/godrop/services/auth/internal/application/dto"
)

type SessionUseCase interface {
	CreateNewSession(ctx context.Context, userID string, metadataDTO dto.SessionMetadataDTO) (*dto.TokenPairDTO, error)
	RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairDTO, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeSessionWithAccessToken(ctx context.Context, accessToken string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	EnsureAccessTokenValid(ctx context.Context, AccessT string) error
	GetUserSessions(ctx context.Context, userID string) ([]*dto.SessionDTO, error)
}
