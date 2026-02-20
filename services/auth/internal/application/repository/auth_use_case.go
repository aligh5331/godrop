package repository

import (
	"context"

	"github.com/aligh5331/godrop/services/auth/internal/application/dto"
)

type AuthUseCase interface {
	Login(ctx context.Context, inputDTO dto.LoginInputDTO, metadataDTO dto.SessionMetadataDTO) (*dto.LoginDTO, error)
	Register(ctx context.Context, inputDTO dto.RegisterInputDTO, metadataDTO dto.SessionMetadataDTO) (*dto.RegisterDTO, error)
	ChangePassword(ctx context.Context, inputDTO dto.ChangePasswordDTO, metadataDTO dto.SessionMetadataDTO) (*dto.LoginDTO, error)
	UpdateName(ctx context.Context, inputDTO dto.UpdateUserNameDTO) error
	UpdateEmail(ctx context.Context, inputDTO dto.UpdateUserEmailDTO) error
	Delete(ctx context.Context, inputDTO dto.DeleteUserDTO) error
}
