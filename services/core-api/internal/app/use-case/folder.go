package use_case

import (
	"context"
	"strings"
	"time"

	"github.com/aligh5331/godrop/services/core-api/internal/app/dto"
	"github.com/aligh5331/godrop/services/core-api/internal/app/helpers"
	"github.com/aligh5331/godrop/services/core-api/internal/domain"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/entity"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/repo"
)

type FolderUseCase struct {
	repo  repo.CoreRepository
	idGen helpers.IdGenerator
}

func (uc *FolderUseCase) GetFolder(ctx context.Context, folderID, userID string) (*dto.Folder, error) {
	folder, err := uc.repo.GetFolder(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	pID := ""
	if folder.ParentId() != nil {
		pID = *folder.ParentId()
	}
	return &dto.Folder{
		ID:       folder.Id(),
		Name:     folder.Name(),
		UserID:   folder.UserId(),
		ParentID: pID,
	}, nil
}

func (uc *FolderUseCase) CreateFolder(ctx context.Context, in dto.CreateFolderInputs) (*dto.Folder, error) {
	pID := strings.TrimSpace(in.ParentID)
	if pID == "" {
		return nil, domain.ErrEmptyParentID
	}
	fID := uc.idGen.NewID()
	now := time.Now()
	folder, err := entity.NewFolder(fID, in.Name, in.UserID, &pID, now, now)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return &dto.Folder{
		ID:       folder.Id(),
		Name:     folder.Name(),
		UserID:   folder.UserId(),
		ParentID: pID,
	}, nil
}

func (uc *FolderUseCase) MoveFolder(ctx context.Context, in dto.MoveFolderInputs) (*dto.Folder, error) {
	err := validateMoveInputs(in)
	if err != nil {
		return nil, err
	}

	ctx, tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	isDescendant, err := uc.repo.IsDescendant(ctx, in.FolderID, in.NewParentID)
	if err != nil {
		return nil, err
	}
	if isDescendant {
		return nil, domain.ErrFoldersCantCreateCycles
	}

	folder, err := uc.repo.MoveFolder(ctx, in.FolderID, in.NewParentID, in.UserID)

	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	pID := ""
	if folder.ParentId() != nil {
		pID = *folder.ParentId()
	}
	return &dto.Folder{
		ID:       folder.Id(),
		Name:     folder.Name(),
		UserID:   folder.UserId(),
		ParentID: pID,
	}, nil
}

func validateMoveInputs(in dto.MoveFolderInputs) error {
	pid := strings.TrimSpace(in.NewParentID)
	fid := strings.TrimSpace(in.FolderID)
	if pid == "" {
		return domain.ErrEmptyParentID
	}
	if fid == "" {
		return domain.ErrEmptyID
	}
	if fid == pid {
		return domain.ErrFolderCantBeItsOwnParent
	}
	return nil
}

func (uc *FolderUseCase) UpdateFolderName(ctx context.Context, in dto.UpdateFolderName) (*dto.Folder, error) {
	folder, err := uc.repo.UpdateFolderName(ctx, in.FolderID, in.NewName, in.UserID)
	if err != nil {
		return nil, err
	}
	pID := ""
	if folder.ParentId() != nil {
		pID = *folder.ParentId()
	}
	return &dto.Folder{
		ID:       folder.Id(),
		Name:     folder.Name(),
		UserID:   folder.UserId(),
		ParentID: pID,
	}, nil
}

func (uc *FolderUseCase) DeleteFolder(ctx context.Context, folderID, userID string) error {
	if err := uc.repo.DeleteFolder(ctx, folderID, userID); err != nil {
		return err
	}
	return nil
}
