package security

import "github.com/google/uuid"

type GoogleUUIDGen struct{}

func (g *GoogleUUIDGen) NewId() string {
	return uuid.Must(uuid.NewV7()).String()
}
