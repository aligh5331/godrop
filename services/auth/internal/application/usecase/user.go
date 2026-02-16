package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"auth/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AuthUseCase struct {
	repo          repository.UserRepository
	hasher        repository.Hasher
	idGen         repository.IdGenerator
	validator     repository.Validator
	emailVerifier repository.EmailVerifier
}

func NewAuthUseCase(
	repo repository.UserRepository,
	hasher repository.Hasher,
	idGen repository.IdGenerator,
	validator repository.Validator,
	emailVerifier repository.EmailVerifier,
) *AuthUseCase {
	return &AuthUseCase{
		repo:          repo,
		hasher:        hasher,
		idGen:         idGen,
		validator:     validator,
		emailVerifier: emailVerifier,
	}
}

func (uc *AuthUseCase) VerifyEmail(ctx context.Context, email string) error {

	email = strings.TrimSpace(email)
	if err := uc.validator.ValidateEmail(email); err != nil {
		return err
	}

	if err := uc.emailVerifier.VerifyEmail(ctx, email); err != nil {
		return err
	}
	return nil
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
			return fmt.Errorf("repo find user by email: %w", err)
		}
	} else {
		return domain.ErrEmailAlreadyExists
	}

	if !uc.emailVerifier.IsEmailVerified(ctx, email) {
		return domain.ErrEmailNotVerified
	}

	hp, hErr := uc.hasher.Hash(password)
	if hErr != nil {
		return fmt.Errorf("hash password: %w", hErr)
	}

	newId := uc.idGen.NewId()
	now := time.Now()
	du, dErr := entity.NewUser(newId, name, email, hp, now, now)
	if dErr != nil {
		return dErr
	}

	if err := uc.repo.Create(ctx, du); err != nil {
		return fmt.Errorf("repo create : %w", err)
	}
	return nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	if err := uc.validator.ValidateEmail(email); err != nil {
		return nil, err
	}

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("repo find user by email: %w", err)
	}

	if !uc.hasher.CheckPasswordHash(password, user.Password()) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {

	//password check
	u, uErr := uc.checkPassword(ctx, userID, oldPassword)
	if uErr != nil {
		return uErr
	}

	//changing password
	if err := uc.validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	hp, hErr := uc.hasher.Hash(newPassword)
	if hErr != nil {
		return fmt.Errorf("hash password: %w", hErr)
	}

	if err := u.ChangePassword(hp); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) UpdateName(ctx context.Context, name, userID string) error {
	u, uErr := uc.repo.FindById(ctx, userID)
	if uErr != nil {
		if errors.Is(uErr, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("repo find user by id: %w", uErr)
	}

	if err := u.ChangeName(name); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}
	return nil
}

func (uc *AuthUseCase) UpdateEmail(ctx context.Context, userID, password, newEmail string) error {

	newEmail = strings.TrimSpace(newEmail)
	if err := uc.validator.ValidateEmail(newEmail); err != nil {
		return err
	}

	u, uErr := uc.checkPassword(ctx, userID, password)
	if uErr != nil {
		return uErr
	}

	if !uc.emailVerifier.IsEmailVerified(ctx, newEmail) {
		return domain.ErrEmailNotVerified
	}

	if err := u.ChangeEmail(newEmail); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}

// password check
func (uc *AuthUseCase) checkPassword(ctx context.Context, userID, password string) (*entity.User, error) {
	u, uErr := uc.repo.FindById(ctx, userID)
	if uErr != nil {
		if errors.Is(uErr, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("repo find user by id: %w", uErr)
	}

	if !uc.hasher.CheckPasswordHash(password, u.Password()) {
		return nil, domain.ErrInvalidCredentials
	}
	return u, nil
}
