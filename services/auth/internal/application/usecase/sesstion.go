package usecase

import (
	"auth/internal/application/dto"
	"auth/internal/domain"
	"auth/internal/domain/entity"
	"auth/internal/domain/helpers"
	"auth/internal/domain/repository"
	"context"
	"errors"
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

	hRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var refreshTokenE *entity.RefreshToken
	//cache check
	cachedData, err := uc.cache.Get(ctx, "rt:"+string(hRefreshToken))

	if err == nil && cachedData != "" {
		refreshTokenE = entity.Unserialize(entity.SerializeRefreshTokenE(cachedData))
	}

	if refreshTokenE == nil {
		refreshTokenE, err = uc.repo.GetRefreshTokenEntityByRefreshToken(ctx, hRefreshToken)
		if err != nil {
			return nil, err
		}
	}

	//rt validation
	if err = refreshTokenE.EnsureValid(now); err != nil {
		if errors.Is(err, domain.ErrReUsedToken) {
			_ = uc.RevokeSession(ctx, refreshTokenE.SessionID())
			return nil, err
		}
		return nil, err
	}

	//generating new tokens
	session, err := uc.repo.GetSessionByID(ctx, refreshTokenE.SessionID())
	if err != nil {
		return nil, err
	}

	metadataDTO := dto.SessionMetadataDTO{IP: session.IP(), ClientAgent: session.UserAgent()}
	newAccessT, err := uc.tokenGen.GenerateAccessToken(
		session.UserID(),
		metadataDTO,
		uc.sDuration,
	)
	if err != nil {
		return nil, err
	}

	newHAccessT, err := uc.hasher.Hash(newAccessT)
	if err != nil {
		return nil, err
	}

	newRefreshT, err := uc.tokenGen.GenerateRefreshToken(session.UserID(), session.ID(), uc.rtDuration)
	if err != nil {
		return nil, err
	}

	newHRefreshT, err := uc.hasher.Hash(newRefreshT)
	if err != nil {
		return nil, err
	}

	ID := uc.idGen.NewId()
	newRefreshTE, err := entity.NewRefreshToken(ID, session.ID(), newHRefreshT, now, uc.rtDuration)
	if err != nil {
		return nil, err
	}
	//remove cache
	if err = uc.cache.Delete(ctx, "at:"+string(session.Token())); err != nil {
		return nil, err
	}
	if err = uc.cache.Delete(ctx, "rt:"+string(hRefreshToken)); err != nil {
		return nil, err
	}

	//update repo
	session.SetToken(newHAccessT, now, uc.sDuration)

	if err = uc.repo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}
	if err = uc.repo.RevokeRefreshToken(ctx, refreshTokenE.ID()); err != nil {
		return nil, err
	}

	if err = uc.repo.CreateRefreshToken(ctx, newRefreshTE); err != nil {
		return nil, err
	}

	val := newRefreshTE.Serialize()
	if err = uc.cache.Set(ctx, "at:"+string(newHAccessT), metadataDTO, uc.sDuration); err != nil {
		return nil, err
	}
	if err = uc.cache.Set(ctx, "rt:"+string(newHRefreshT), string(val), uc.rtDuration); err != nil {
		return nil, err
	}

	return &dto.TokenPairDTO{
		RefreshToken: newRefreshT,
		AccessToken:  newAccessT,
	}, nil
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

func (uc *SessionUseCase) EnsureAccessTokenValid(ctx context.Context, AccessT string) error {

	hAccessT, err := uc.hasher.Hash(AccessT)
	if err != nil {
		return err
	}
	if ok, _ := uc.cache.Exists(ctx, "at:"+string(hAccessT)); ok {
		return nil
	}

	session, err := uc.repo.GetSessionByAccessToken(ctx, hAccessT)
	if err != nil {
		return err
	}

	now := time.Now()
	if !session.IsValid(now) {
		return domain.ErrInvalidSession
	}
	session.Use(now)
	_ = uc.repo.UpdateSession(ctx, session)
	return nil
}
