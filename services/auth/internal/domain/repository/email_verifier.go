package repository

import "context"

type EmailVerifier interface {
	IsVerified(ctx context.Context, email string) bool
	Verify(ctx context.Context, email string) error
}
