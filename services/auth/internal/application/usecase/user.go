package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"auth/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type AuthUseCase struct {
	repo      repository.UserRepository
	hasher    repository.Hasher
	idGen     repository.IdGenerator
	validator repository.Validator
}

func NewAuthUseCase(
	repo repository.UserRepository,
	hasher repository.Hasher,
	idGen repository.IdGenerator,
	validator repository.Validator,
) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		hasher:    hasher,
		idGen:     idGen,
		validator: validator,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, name, password string) error {
	if err := uc.validator.ValidatePassword(password); err != nil {
		return err
	}

	if err := uc.validator.ValidateEmail(email); err != nil {
		return err
	}

	if _, err := uc.repo.FindByEmail(ctx, email); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return fmt.Errorf("find user by email: %w", err)
		}
	} else {
		return domain.ErrEmailAlreadyExists
	}

	hp, hErr := uc.hasher.Hash(password)
	if hErr != nil {
		return fmt.Errorf("hash password: %w", hErr)
	}

	newId := uc.idGen.NewId()
	now := time.Now()
	du, dErr := entity.NewUser(newId, name, email, hp, true, now, now)
	if dErr != nil {
		return dErr
	}

	rErr := uc.repo.Create(ctx, du)
	if rErr != nil {
		return fmt.Errorf("repo create : %w", rErr)
	}
	return nil
}

func (uc *AuthUseCase) login(ctx context.Context, email, password string) (*entity.User, error) {
	if err := uc.validator.ValidateEmail(email); err != nil {
		return nil, err
	}

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	if !uc.hasher.CheckPasswordHash(password, user.Password()) {
		return nil, domain.ErrWrongPassword
	}

	return user, nil
}
