package use_case

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/service"
)

type UserUseCase struct {
	auth service.AuthService
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*dto.User, error) {
	return uc.auth.GetUser(ctx, id)
}

func (uc *UserUseCase) ChangeUserName(ctx context.Context, id string, newName string) (*dto.User, error) {

	//TODO: Add User return to auth.UpdateName service

	err := uc.auth.UpdateName(ctx, id, newName)
	if err != nil {
		return nil, err
	}
	user, err := uc.auth.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id, pass string) error {
	return uc.auth.DeleteUser(ctx, id, pass)
}
