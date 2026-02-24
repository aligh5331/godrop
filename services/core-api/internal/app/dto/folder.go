package dto

type CreateFolderInputs struct {
	ParentID string `json:"parent_id"`
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
}

type Folder struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	ParentID string `json:"parent_id"`
}

type MoveFolderInputs struct {
	FolderID string `json:"folder_id"`
	ParentID string `json:"parent_id"`
}

type UpdateFolderName struct {
	FolderID string `json:"folder_id"`
	NewName  string `json:"new_name"`
}
