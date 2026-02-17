package dto

type RegisterDTO struct {
	User   *UserDTO
	Tokens *TokenPairDTO
}

type LoginDTO struct {
	User   *UserDTO
	Tokens *TokenPairDTO
}
type RegisterInputDTO struct {
	Email    string
	Password string
	Name     string
}

type LoginInputDTO struct {
	Email    string
	Password string
}

type ChangePasswordDTO struct {
	UserID  string
	OldPass string
	NewPass string
}

type UserDTO struct {
	ID    string
	Name  string
	Email string
}

type UpdateUserEmailDTO struct {
	UserID string
	Pass   string
	Email  string
}

type UpdateUserNameDTO struct {
	UserID string
	Name   string
}

type DeleteUserDTO struct {
	UserID string
	Pass   string
}
