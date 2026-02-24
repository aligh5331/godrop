package domain

import "errors"

var (
	ErrEmptyID       = errors.New("empty ID")
	ErrEmptyUserID   = errors.New("empty user ID")
	ErrEmptyName     = errors.New("empty Name")
	ErrEmptyParentID = errors.New("empty ParentID")
)
