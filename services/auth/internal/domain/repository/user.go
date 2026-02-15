package repository

import (
	"auth/internal/domain/entity"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindById(ctx context.Context, userUUID string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
}
