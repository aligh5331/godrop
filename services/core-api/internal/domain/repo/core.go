package repo

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/domain/entity"
)

type CoreRepository interface {
	//User
	CreateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	//Folder
	GetFolder(ctx context.Context, id string) (*entity.Folder, error)
	CreateFolder(ctx context.Context, folder *entity.Folder) error
	UpdateFolderName(ctx context.Context, folderId, newName string) (*entity.Folder, error)
	MoveFolder(ctx context.Context, folderId, newParentId string) (*entity.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
	//Transactoin
	BeginTx(ctx context.Context) (context.Context, Transaction, error)
}
