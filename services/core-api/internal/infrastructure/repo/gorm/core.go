package gorm

import (
	"context"
	"time"

	"github.com/aligh5331/godrop/services/core-api/internal/domain"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/entity"
	"github.com/aligh5331/godrop/services/core-api/internal/domain/repo"
	"gorm.io/gorm"
)

type CoreRepository struct {
	db *gorm.DB
}

func (r *CoreRepository) CreateUser(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx, r.db)
	u := ToGromUser(user)
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *CoreRepository) DeleteUser(ctx context.Context, id string) error {
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Delete(&entity.User{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CoreRepository) GetFolder(ctx context.Context, folderID, userID string) (*entity.Folder, error) {
	folder := &GromFolder{}
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Model(&GromFolder{}).Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return folder.ToEntity(), nil
}

func (r *CoreRepository) CreateFolder(ctx context.Context, folder *entity.Folder) error {
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Create(ToGromFolder(folder)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CoreRepository) UpdateFolderName(ctx context.Context, folderId, newName, userID string) (*entity.Folder, error) {
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Model(&GromFolder{}).Where("id = ? AND user_id = ?", folderId, userID).Update("name", newName).Error
	if err != nil {
		return nil, err
	}

	folder, err := r.GetFolder(ctx, folderId, userID)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *CoreRepository) MoveFolder(ctx context.Context, folderID, newParentID, userID string) (*entity.Folder, error) {
	if folderID == newParentID {
		return nil, domain.ErrFolderCantBeItsOwnParent
	}

	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&GromFolder{}).
		Where("id = ? AND user_id = ?", folderID, userID).
		Update("parent_id", newParentID).Error
	if err != nil {
		return nil, err
	}

	folder := &GromFolder{}
	err = db.WithContext(ctx).
		Model(&GromFolder{}).
		Where("id = ? AND user_id = ?", folderID, userID).
		First(&folder).Error
	if err != nil {
		return nil, err
	}

	return folder.ToEntity(), nil
}

func (r *CoreRepository) DeleteFolder(ctx context.Context, folderID, userID string) error {
	now := time.Now()
	query := `
        WITH RECURSIVE folder_branch AS (
            SELECT id FROM folders WHERE id = ? AND user_id = ? AND deleted_at IS NULL
            UNION ALL
            SELECT f.id FROM folders f INNER JOIN folder_branch fb ON f.parent_id = fb.id
            WHERE f.deleted_at IS NULL
        )
        UPDATE folders 
        SET deleted_at = ?, updated_at = ?
        WHERE id IN (SELECT id FROM folder_branch)`

	return GetDB(ctx, r.db).WithContext(ctx).Exec(query, folderID, userID, now, now).Error
}

// IsDescendant Check if childID is a descendant of the parentID
func (r *CoreRepository) IsDescendant(ctx context.Context, parentID, childID string) (bool, error) {
	var exists bool

	checkCycleQuery := `
    WITH RECURSIVE descendants AS (
        SELECT id FROM folders WHERE id = ?
        UNION ALL
        SELECT f.id FROM folders f INNER JOIN descendants d ON f.parent_id = d.id
    )
    SELECT EXISTS(SELECT 1 FROM descendants WHERE id = ?)`

	db := GetDB(ctx, r.db)

	err := db.Raw(checkCycleQuery, parentID, childID).Scan(&exists).Error

	return exists, err
}

type Transaction struct {
	tx *gorm.DB
}

type GromUser struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *GromUser) ToEntity() *entity.User {
	return entity.NewUser(u.ID, u.CreatedAt)
}

func ToGromUser(u *entity.User) *GromUser {
	return &GromUser{
		ID:        u.ID(),
		CreatedAt: u.CreatedAt(),
	}
}

type GromFolder struct {
	ID        string  `gorm:"primarykey"`
	UserID    string  `gorm:"not null"`
	ParentID  *string `gorm:"column:parent_id"`
	Name      string  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Parent   *GromFolder  `gorm:"foreignKey:ParentID"`
	Children []GromFolder `gorm:"foreignKey:ParentID"`
}

func ToGromFolder(folder *entity.Folder) *GromFolder {
	return &GromFolder{
		ID:        folder.Id(),
		Name:      folder.Name(),
		ParentID:  folder.ParentId(),
		UserID:    folder.UserId(),
		CreatedAt: folder.CreatedAt(),
		UpdatedAt: folder.UpdatedAt(),
	}
}

func (f *GromFolder) ToEntity() *entity.Folder {
	folder, _ := entity.NewFolder(f.ID, f.Name, f.UserID, f.ParentID, f.CreatedAt, f.UpdatedAt)
	return folder
}

func (f *GromFolder) TableName() string {
	return "folders"
}

func (u *GromUser) TableName() string {
	return "users"
}

func (t *Transaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *Transaction) Tx() *gorm.DB {
	return t.tx
}

func (r *CoreRepository) GetFolderWithChildren(ctx context.Context, folderID string) (*GromFolder, error) {
	db := GetDB(ctx, r.db)
	var folder GromFolder

	// Preload triggers a second query: SELECT * FROM folders WHERE parent_id = ?
	err := db.Preload("Children").First(&folder, "id = ?", folderID).Error

	if err != nil {
		return nil, err
	}

	return &folder, nil
}
func (r *CoreRepository) BeginTx(ctx context.Context) (context.Context, repo.Transaction, error) {
	if _, ok := ctx.Value(repo.TxKey{}).(*gorm.DB); ok {
		// A transaction already exists. Return the context as-is
		// and a NO-OP transaction so the caller doesn't break the parent TX.
		return ctx, &noOpTx{}, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	newCtx := context.WithValue(ctx, repo.TxKey{}, tx)
	return newCtx, &Transaction{tx: tx}, nil
}

//func (r *CoreRepository) CreateFolderWithDefaultName(ctx context.Context, id, userID string, parentID *string) (*GromFolder, error) {
//
//	ctx, tx, err := r.BeginTx(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer tx.Rollback()
//	db := GetDB(ctx, r.db)
//
//	baseName := "New Folder"
//	var finalName string
//
//	var existingNames []string
//
//	// 1. Fetch all names starting with "New Folder" under this parent
//	query := db.WithContext(ctx).Model(&GromFolder{}).Where("user_id = ? AND name LIKE ?", userID, baseName+"%")
//	if parentID == nil {
//		query = query.Where("parent_id IS NULL")
//	} else {
//		query = query.Where("parent_id = ?", *parentID)
//	}
//
//	query.Pluck("name", &existingNames)
//
//	// 2. Determine the next available name
//	finalName = findNextAvailableName(baseName, existingNames)
//
//	// 3. Create the record
//	newFolder := &GromFolder{
//		ID:       id,
//		Name:     finalName,
//		UserID:   userID,
//		ParentID: parentID,
//	}
//
//	if err = db.WithContext(ctx).Create(newFolder).Error; err != nil {
//		return nil, err
//	}
//	if err = tx.Commit(); err != nil {
//		return nil, err
//	}
//
//	return &GromFolder{
//		ID:       id,
//		Name:     finalName,
//		UserID:   userID,
//		ParentID: parentID,
//	}, nil
//}
//
//func findNextAvailableName(base string, existing []string) string {
//	nameMap := make(map[string]bool)
//	for _, n := range existing {
//		nameMap[n] = true
//	}
//
//	if !nameMap[base] {
//		return base
//	}
//
//	for i := 1; ; i++ {
//		candidate := fmt.Sprintf("%s (%d)", base, i)
//		if !nameMap[candidate] {
//			return candidate
//		}
//	}
//}
