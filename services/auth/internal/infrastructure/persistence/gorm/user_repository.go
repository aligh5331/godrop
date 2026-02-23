package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/domain"
	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"
	"github.com/aligh5331/godrop/services/auth/internal/domain/helpers"
	"github.com/aligh5331/godrop/services/auth/internal/domain/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type URTransaction struct {
	tx *gorm.DB
}

func (t *URTransaction) Tx() *gorm.DB {
	return t.tx
}
func (t *URTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *URTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (u *UserRepository) BeginTx(ctx context.Context) (context.Context, repository.Transaction, error) {
	if _, ok := ctx.Value(helpers.TxKey{}).(*gorm.DB); ok {
		// A transaction already exists. Return the context as-is
		// and a NO-OP transaction so the caller doesn't break the parent TX.
		return ctx, &noOpTx{}, nil
	}

	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	newCtx := context.WithValue(ctx, helpers.TxKey{}, tx)
	return newCtx, &SRTransaction{tx: tx}, nil
}

func (u *UserRepository) Create(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx, u.db)
	gu := toGormUser(user)
	if err := db.WithContext(ctx).Create(gu).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) FindById(ctx context.Context, userUUID string) (*entity.User, error) {
	db := GetDB(ctx, u.db)
	gu := &gUser{}
	if err := db.WithContext(ctx).Where("id = ?", userUUID).First(gu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("gorm find by id: %w", err)
	}
	user, dErr := toDomainUser(gu)
	if dErr != nil {
		return nil, dErr
	}
	return user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := GetDB(ctx, u.db)

	var gu = &gUser{}
	if err := db.WithContext(ctx).Where("email = ?", email).First(gu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("gorm find by email: %w", err)
	}
	user, dErr := toDomainUser(gu)
	if dErr != nil {
		return nil, dErr
	}
	return user, nil
}

func (u *UserRepository) Update(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx, u.db)

	res := db.WithContext(ctx).
		Model(&gUser{}).
		Where("id = ?", user.Id()).
		Updates(map[string]interface{}{
			"name":       user.Name(),
			"email":      user.Email(),
			"password":   string(user.Password()),
			"updated_at": user.UpdatedAt(),
		})

	if res.Error != nil {
		return fmt.Errorf("gorm update user: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (u *UserRepository) Delete(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx, u.db)

	gu := toGormUser(user)
	if err := db.WithContext(ctx).Delete(gu).Error; err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}

type gUser struct {
	ID        string `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(255);not null"` // bcrypt outputs 60 chars
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (gUser) TableName() string {
	return "users"
}

func toGormUser(u *entity.User) *gUser {
	return &gUser{
		ID:        u.Id(),
		Name:      u.Name(),
		Email:     u.Email(),
		Password:  string(u.Password()),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func toDomainUser(g *gUser) (*entity.User, error) {
	return entity.NewUser(
		g.ID,
		g.Name,
		g.Email,
		entity.HashedPassword(g.Password),
		g.CreatedAt,
		g.UpdatedAt,
	)
}
