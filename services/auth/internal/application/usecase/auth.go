package usecase

import (
	"auth/internal/application/dto"
	ar "auth/internal/application/repository"
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"auth/internal/domain/helpers"
	dr "auth/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AuthUseCase struct {
	sessionUc     ar.SessionUseCase
	repo          dr.UserRepository
	hasher        helpers.PasswordHasher
	idGen         helpers.IdGenerator
	validator     helpers.Validator
	emailVerifier helpers.EmailVerifier
}

func NewAuthUseCase(
	sessionUc ar.SessionUseCase,
	repo dr.UserRepository,
	hasher helpers.PasswordHasher,
	idGen helpers.IdGenerator,
	validator helpers.Validator,
	emailVerifier helpers.EmailVerifier,
) ar.AuthUseCase {
	return &AuthUseCase{
		sessionUc:     sessionUc,
		repo:          repo,
		hasher:        hasher,
		idGen:         idGen,
		validator:     validator,
		emailVerifier: emailVerifier,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, inputDTO dto.LoginInputDTO, metadataDTO dto.SessionMetadataDTO) (*dto.LoginDTO, error) {
	if err := uc.validator.ValidateEmail(inputDTO.Email); err != nil {
		return nil, err
	}

	user, err := uc.repo.FindByEmail(ctx, inputDTO.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("repo find user by email: %w", err)
	}

	if !uc.hasher.CheckPasswordHash(inputDTO.Password, user.Password()) {
		return nil, domain.ErrInvalidCredentials
	}

	tokens, err := uc.sessionUc.CreateNewSession(ctx, user.Id(), metadataDTO)
	if err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	userDTO := &dto.UserDTO{
		ID:    user.Id(),
		Name:  user.Name(),
		Email: user.Email(),
	}

	loginDTO := &dto.LoginDTO{
		User:   userDTO,
		Tokens: tokens,
	}

	return loginDTO, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, inputDTO dto.RegisterInputDTO, metadataDTO dto.SessionMetadataDTO) (*dto.RegisterDTO, error) {
	if err := uc.validator.ValidatePassword(inputDTO.Password); err != nil {
		return nil, err
	}

	if err := uc.validator.ValidateEmail(inputDTO.Email); err != nil {
		return nil, err
	}

	if _, err := uc.repo.FindByEmail(ctx, inputDTO.Email); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("repo find user by email: %w", err)
		}
	} else {
		return nil, domain.ErrEmailAlreadyExists
	}

	if !uc.emailVerifier.IsEmailVerified(ctx, inputDTO.Email) {
		return nil, domain.ErrEmailNotVerified
	}

	hp, err := uc.hasher.Hash(inputDTO.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newId := uc.idGen.NewId()
	now := time.Now()
	user, err := entity.NewUser(newId, inputDTO.Name, inputDTO.Email, hp, now, now)
	if err != nil {
		return nil, err
	}

	if err = uc.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("repo create : %w", err)
	}

	tokens, err := uc.sessionUc.CreateNewSession(ctx, user.Id(), metadataDTO)
	if err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	userDTO := &dto.UserDTO{
		ID:    user.Id(),
		Name:  user.Name(),
		Email: user.Email(),
	}

	registerDTO := &dto.RegisterDTO{
		User:   userDTO,
		Tokens: tokens,
	}

	return registerDTO, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, inputDTO dto.ChangePasswordDTO, metadataDTO dto.SessionMetadataDTO) (*dto.LoginDTO, error) {
	//password check
	user, err := uc.checkPassword(ctx, inputDTO.UserID, inputDTO.OldPass)
	if err != nil {
		return nil, err
	}

	//changing password
	if err = uc.validator.ValidatePassword(inputDTO.NewPass); err != nil {
		return nil, err
	}

	hp, err := uc.hasher.Hash(inputDTO.NewPass)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	if err = user.ChangePassword(hp, time.Now()); err != nil {
		return nil, err
	}

	if err = uc.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("repo update: %w", err)
	}

	if err = uc.sessionUc.RevokeAllUserSessions(ctx, user.Id()); err != nil {
		return nil, err
	}

	tokens, err := uc.sessionUc.CreateNewSession(ctx, user.Id(), metadataDTO)
	if err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	userDTO := &dto.UserDTO{
		ID:    user.Id(),
		Name:  user.Name(),
		Email: user.Email(),
	}

	loginDTO := &dto.LoginDTO{
		User:   userDTO,
		Tokens: tokens,
	}

	return loginDTO, nil
}

func (uc *AuthUseCase) UpdateName(ctx context.Context, inputDTO dto.UpdateUserNameDTO) error {
	u, err := uc.repo.FindById(ctx, inputDTO.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("repo find user by id: %w", err)
	}

	if err = u.ChangeName(inputDTO.Name, time.Now()); err != nil {
		return err
	}

	if err = uc.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}
	return nil
}

func (uc *AuthUseCase) UpdateEmail(ctx context.Context, inputDTO dto.UpdateUserEmailDTO) error {
	newEmail := strings.TrimSpace(inputDTO.Email)
	if err := uc.validator.ValidateEmail(newEmail); err != nil {
		return err
	}

	u, err := uc.checkPassword(ctx, inputDTO.UserID, inputDTO.Pass)
	if err != nil {
		return err
	}

	if !uc.emailVerifier.IsEmailVerified(ctx, newEmail) {
		return domain.ErrEmailNotVerified
	}

	if err = u.ChangeEmail(newEmail, time.Now()); err != nil {
		return err
	}

	if err = uc.repo.Update(ctx, u); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) Delete(ctx context.Context, inputDTO dto.DeleteUserDTO) error {
	user, err := uc.checkPassword(ctx, inputDTO.UserID, inputDTO.Pass)
	if err != nil {
		return err
	}

	if err = uc.sessionUc.RevokeAllUserSessions(ctx, user.Id()); err != nil {
		return err
	}
	if err = uc.repo.Delete(ctx, user); err != nil {
		return err
	}
	return nil
}

func (uc *AuthUseCase) checkPassword(ctx context.Context, userID, password string) (*entity.User, error) {
	u, err := uc.repo.FindById(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("repo find user by id: %w", err)
	}

	if !uc.hasher.CheckPasswordHash(password, u.Password()) {
		return nil, domain.ErrInvalidCredentials
	}
	return u, nil
}
