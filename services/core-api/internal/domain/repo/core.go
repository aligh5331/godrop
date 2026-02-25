package repo

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/domain/entity"
)

type CoreRepository interface {
	//User
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	//Folder
	GetFolder(ctx context.Context, id string) (*entity.Folder, error)
	SaveFolder(ctx context.Context, folder *entity.Folder) (*entity.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
	//Transactoin
	BeginTx(ctx context.Context) (context.Context, Transaction, error)
}
