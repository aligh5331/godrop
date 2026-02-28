package use_case

import (
	"context"
	"time"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
	"github.com/aligh5331/godrop/services/core-api/internal/app/helpers"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/entity"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/repo"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/service"
)

type AuthUseCse struct {
	repo      repo.CoreRepository
	auth      service.AuthService
	validator helpers.Validator
	idGen     helpers.IdGenerator
}

func (uc *AuthUseCse) Login(ctx context.Context, in dto.Login) (*dto.TokenPairs, error) {
	if err := uc.validator.Email(in.Email); err != nil {
		return nil, err
	}

	tokens, err := uc.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUseCse) Register(ctx context.Context, in dto.Register) (*dto.TokenPairs, error) {
	now := time.Now()
	tokens, err := uc.auth.Register(ctx, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	var success bool
	defer func() {
		if !success {
			_ = uc.auth.DeleteUser(ctx, tokens.User.ID, in.Password)
		}
	}()

	ctx, tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user := entity.NewUser(tokens.User.ID, now)
	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	rootFolder, err := entity.NewRootFolder(uc.idGen.NewID(), tokens.User.ID, now)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.CreateFolder(ctx, rootFolder); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	success = true
	return tokens, nil
}

func (uc *AuthUseCse) ChangePassword(ctx context.Context, in dto.ChangePasswordInput) (*dto.TokenPairs, error) {
	tokens, err := uc.auth.ChangePassword(ctx, in.UserID, in.NewPassword, in.OldPassword)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUseCse) ChangeEmail(ctx context.Context, in dto.ChangeEmailInput) (*dto.TokenPairs, error) {
	tokens, err := uc.auth.UpdateEmail(ctx, in.UserID, in.NewEmail, in.Password)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUseCse) RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairs, error) {
	tokens, err := uc.auth.RefreshSession(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUseCse) Logout(ctx context.Context, accessToken string) error {
	err := uc.auth.Logout(ctx, accessToken)
	if err != nil {
		return err
	}
	return nil
}

func (uc *AuthUseCse) RevokeSession(ctx context.Context, sessionID string) error {
	err := uc.auth.RevokeSession(ctx, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *AuthUseCse) RevokeAllSessions(ctx context.Context, userID string) error {
	err := uc.auth.RevokeAllSessions(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *AuthUseCse) GetUserSessions(ctx context.Context, userID string) ([]*dto.Session, error) {
	out, err := uc.auth.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return out, nil
}
