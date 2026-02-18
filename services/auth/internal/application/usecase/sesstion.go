package usecase

import (
	"auth/internal/application/dto"
	"auth/internal/domain/entity"
	"auth/internal/domain/helpers"
	"auth/internal/domain/repository"
	"context"
	"time"
)

type SessionUseCase struct {
	sDuration  time.Duration
	rtDuration time.Duration
	repo       repository.SessionRepository
	cache      repository.CacheRepository
	idGen      helpers.IdGenerator
	hasher     helpers.TokenHasher
	tokenGen   helpers.TokenGenerator
	validator  helpers.Validator
}

func (uc *SessionUseCase) CreateNewSession(
	ctx context.Context,
	userID string,
	metadataDTO dto.SessionMetadataDTO,
) (*dto.TokenPairDTO, error) {
	if err := uc.validator.ValidateIP(metadataDTO.IP); err != nil {
		return nil, err
	}

	now := time.Now()

	ID := uc.idGen.NewId()

	//generate session token
	sessionToken, err := uc.tokenGen.GenerateAccessToken(userID, metadataDTO, uc.sDuration)
	if err != nil {
		return nil, err
	}

	//hash it
	hSessionToken, err := uc.hasher.Hash(sessionToken)
	if err != nil {
		return nil, err
	}

	//generate refresh token
	refreshToken, err := uc.tokenGen.GenerateRefreshToken(userID, ID, uc.rtDuration)
	if err != nil {
		return nil, err
	}

	//hash it
	hRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	sessionE, err := entity.NewSession(
		ID,
		metadataDTO.ClientAgent,
		metadataDTO.IP,
		userID,
		hSessionToken,
		now, uc.sDuration,
	)
	if err != nil {
		return nil, err
	}
	refreshTokenE, err := entity.NewRefreshToken(ID, ID, hRefreshToken, now, uc.rtDuration)
	if err != nil {
		return nil, err
	}

	if err = uc.repo.CreateSession(ctx, sessionE); err != nil {
		return nil, err
	}
	if err = uc.repo.CreateRefreshToken(ctx, refreshTokenE); err != nil {
		return nil, err
	}

	val := refreshTokenE.Serialize()
	if err = uc.cache.Set(ctx, "at:"+string(hSessionToken), metadataDTO, uc.sDuration); err != nil {
		return nil, err
	}
	if err = uc.cache.Set(ctx, "rt:"+string(hRefreshToken), string(val), uc.rtDuration); err != nil {
		return nil, err
	}

	return &dto.TokenPairDTO{
		RefreshToken: refreshToken,
		AccessToken:  sessionToken,
	}, nil
}

func (uc *SessionUseCase) RefreshSession(ctx context.Context, refreshToken string) (*dto.TokenPairDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (uc *SessionUseCase) RevokeSession(ctx context.Context, sessionID string) error {
	//TODO implement me
	panic("implement me")
}

func (uc *SessionUseCase) RevokeAllUserSessions(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (uc *SessionUseCase) GetUserSessions(ctx context.Context, userID string) ([]*dto.SessionDTO, error) {
	//TODO implement me
	panic("implement me")
}
