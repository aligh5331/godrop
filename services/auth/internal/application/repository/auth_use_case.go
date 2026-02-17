package repository

import (
	"auth/internal/application/dto"
	"context"
)

type AuthUseCase interface {
	Login(ctx context.Context, inputDTO dto.LoginInputDTO) (*dto.LoginDTO, error)
	Register(ctx context.Context, inputDTO dto.RegisterInputDTO) (*dto.RegisterDTO, error)
	ChangePassword(ctx context.Context, inputDTO dto.ChangePasswordDTO) (*dto.LoginDTO, error)
	UpdateName(ctx context.Context, inputDTO dto.UpdateUserNameDTO) error
	UpdateEmail(ctx context.Context, inputDTO dto.UpdateUserEmailDTO) error
	Delete(ctx context.Context, inputDTO dto.DeleteUserDTO) error
}
