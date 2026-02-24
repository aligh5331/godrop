package entity

import (
	"strings"

	"github.com/aligh5331/godrop/services/core-api/internal/domain"
)

type Folder struct {
	id       string
	name     string
	userId   string
	parentId *string
}

func NewFolder(id string, name string, userId string, parentId string) (*Folder, error) {
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
	parentId = strings.TrimSpace(parentId)
	if parentId == "" {
		return nil, domain.ErrEmptyParentID
	}
	return &Folder{
		id:       id,
		name:     name,
		userId:   userId,
		parentId: &parentId,
	}, nil
}

func NewRootFolder(id string, userId string) (*Folder, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.ErrEmptyID
	}
	name := "root"
	return &Folder{
		id:     id,
		name:   name,
		userId: userId,
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
