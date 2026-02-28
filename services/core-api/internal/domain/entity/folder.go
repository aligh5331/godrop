package entity

import (
	"strings"
	"time"

	"github.com/aligh5331/godrop/services/core-api/internal/domain"
)

type Folder struct {
	id        string
	name      string
	userId    string
	parentId  *string
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func NewFolder(id, name, userId string, parentId *string, createdAt, updatedAt time.Time) (*Folder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.ErrEmptyID
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	userId = strings.TrimSpace(userId)
	if userId == "" {
		return nil, domain.ErrEmptyUserID
	}
	return &Folder{
		id:        id,
		name:      name,
		userId:    userId,
		parentId:  parentId,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func NewRootFolder(id string, userId string, now time.Time) (*Folder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.ErrEmptyID
	}
	name := "root"
	return &Folder{
		id:        id,
		name:      name,
		userId:    userId,
		parentId:  nil,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (f *Folder) Id() string {
	return f.id
}
func (f *Folder) Name() string {
	return f.name
}
func (f *Folder) UserId() string {
	return f.userId
}
func (f *Folder) ParentId() *string {
	return f.parentId
}
func (f *Folder) CreatedAt() time.Time {
	return f.createdAt
}
func (f *Folder) UpdatedAt() time.Time {
	return f.updatedAt
}
