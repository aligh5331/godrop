package interfaces

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
)

type AuthUseCase interface {
	Login(ctx context.Context, in dto.Login) (*dto.TokenPairs, error)
	Register(ctx context.Context, in dto.Register) (*dto.TokenPairs, error)
	ChangePassword(ctx context.Context, input dto.ChangePasswordInput) (*dto.TokenPairs, error)
	ChangeEmail(ctx context.Context, in dto.ChangeEmailInput) (*dto.TokenPairs, error)
	RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairs, error)
	Logout(ctx context.Context, accessToken string) error
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllSessions(ctx context.Context, userID string) error
	GetUserSessions(ctx context.Context, userID string) ([]*dto.Session, error)
}
