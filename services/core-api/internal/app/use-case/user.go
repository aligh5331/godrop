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
	return uc.auth.UpdateName(ctx, id, newName)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.auth.DeleteUser(ctx, id)
}
