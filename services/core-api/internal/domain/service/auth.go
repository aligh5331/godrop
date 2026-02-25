package service

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*dto.TokenPairs, error)
	Register(ctx context.Context, name, email, password string) (*dto.TokenPairs, error)
	ChangePassword(ctx context.Context, userID, newPass, OldPass string) (*dto.TokenPairs, error)
	UpdateName(ctx context.Context, userID, newName string) (*dto.User, error)
	UpdateEmail(ctx context.Context, userID, newEmail, pass string) (*dto.TokenPairs, error)
	DeleteUser(ctx context.Context, userID string) error
	RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairs, error)
	Logout(ctx context.Context, accessToken string) error
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllSessions(ctx context.Context, userID string) error
	CheckSessionIsValid(ctx context.Context, accessToken string) bool
	GetUserSessions(ctx context.Context, userID string) ([]*dto.Session, error)
}
