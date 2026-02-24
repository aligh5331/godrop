package interfaces

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
)

type UserUceCase interface {
	GetUser(ctx context.Context, id string) (*dto.User, error)
	ChangeUserName(ctx context.Context, id string, name string) (*dto.User, error)
	DeleteUser(ctx context.Context, id string) (*dto.User, error)
}
