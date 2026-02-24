package dto

type Register struct {
	Name     string
	Email    string
	Password string
}

type Login struct {
	Email    string
	Password string
}

type TokenPairs struct {
	User         User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangeEmailInput struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	NewEmail string `json:"new_email"`
}

type Session struct {
	ID        string `json:"id"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type ChangePasswordInput struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
