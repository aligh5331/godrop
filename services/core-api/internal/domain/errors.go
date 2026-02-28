package domain

import "errors"

var (
	ErrEmptyID       = errors.New("empty ID")
	ErrEmptyUserID   = errors.New("empty user ID")
	ErrEmptyName     = errors.New("empty Name")
	ErrEmptyParentID = errors.New("empty NewParentID")
)

var (
	ErrFolderCantBeItsOwnParent = errors.New("folder can't be its own parent")
	ErrFoldersCantCreateCycles  = errors.New("cannot move folder into its own descendant (avoids cycles)")
)
