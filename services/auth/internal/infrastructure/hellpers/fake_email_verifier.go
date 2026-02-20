package hellpers

//TODO make a real one

import (
	"context"
)

type FakeEmailVerifier struct{}

func (f *FakeEmailVerifier) IsEmailVerified(ctx context.Context, email string) bool {
	return true
}

func (f *FakeEmailVerifier) VerifyEmail(ctx context.Context, email string) error {
	return nil
}
