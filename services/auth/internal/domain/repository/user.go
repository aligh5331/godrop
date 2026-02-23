package repository

import (
	"context"

	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindById(ctx context.Context, userUUID string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error

	//Transaction
	BeginTx(ctx context.Context) (context.Context, Transaction, error)
}
