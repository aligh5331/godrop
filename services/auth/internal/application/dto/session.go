package dto

type TokenPairDTO struct {
	AccessToken  string
	RefreshToken string
}

type SessionMetadataDTO struct {
	IP          string
	ClientAgent string
}

type SessionDTO struct {
	SessionId   string
	AccessToken string
	Metadata    SessionMetadataDTO
}
