package repository

import "context"

type EmailVerifier interface {
	IsEmailVerified(ctx context.Context, email string) bool
	VerifyEmail(ctx context.Context, email string) error
}
