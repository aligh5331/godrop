package interfaces

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
)

type FolderUseCase interface {
	GetFolder(ctx context.Context, id string) (*dto.Folder, error)
	CreateFolder(ctx context.Context, in dto.CreateFolderInputs) (*dto.Folder, error)
	MoveFolder(ctx context.Context, in dto.MoveFolderInputs) (*dto.Folder, error)
	UpdateFolderName(ctx context.Context, in dto.UpdateFolderName) (*dto.Folder, error)
	DeleteFolder(ctx context.Context, id string) (*dto.Folder, error)
}
