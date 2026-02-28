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
	GetFolder(ctx context.Context, folderID, userID string) (*entity.Folder, error)
	CreateFolder(ctx context.Context, folder *entity.Folder) error
	UpdateFolderName(ctx context.Context, folderId, newName, userID string) (*entity.Folder, error)
	MoveFolder(ctx context.Context, folderId, newParentId, userID string) (*entity.Folder, error)
	DeleteFolder(ctx context.Context, folderID, userID string) error
	IsDescendant(ctx context.Context, parentID, childID string) (bool, error)
	//Transactoin
	BeginTx(ctx context.Context) (context.Context, Transaction, error)
}
