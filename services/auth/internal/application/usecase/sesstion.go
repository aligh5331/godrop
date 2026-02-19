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

	sID := uc.idGen.NewId()

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
	refreshToken, err := uc.tokenGen.GenerateRefreshToken(userID, sID, uc.rtDuration)
	if err != nil {
		return nil, err
	}

	//hash it
	hRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	sessionE, err := entity.NewSession(
		sID,
		metadataDTO.ClientAgent,
		metadataDTO.IP,
		userID,
		hSessionToken,
		now, uc.sDuration,
	)
	rtID := uc.idGen.NewId()
	if err != nil {
		return nil, err
	}
	refreshTokenE, err := entity.NewRefreshToken(rtID, sID, hRefreshToken, now, uc.rtDuration)
	if err != nil {
		return nil, err
	}

	txRepo, tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // no-op if already committed

	if err = txRepo.CreateSession(ctx, sessionE); err != nil {
		return nil, err
	}
	if err = txRepo.CreateRefreshToken(ctx, refreshTokenE); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	val := refreshTokenE.Serialize()
	_ = uc.cache.Set(ctx, "at:"+string(hSessionToken), metadataDTO, uc.sDuration)
	_ = uc.cache.Set(ctx, "rt:"+string(hRefreshToken), string(val), uc.rtDuration)

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

	//update repo
	txRepo, tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // no-op if already committed

	session.SetToken(newHAccessT, now, uc.sDuration)

	if err = txRepo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}
	if err = txRepo.RevokeRefreshToken(ctx, refreshTokenE.ID()); err != nil {
		return nil, err
	}

	if err = txRepo.CreateRefreshToken(ctx, newRefreshTE); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	//remove cache
	_ = uc.cache.Delete(ctx, "at:"+string(session.Token()))
	_ = uc.cache.Delete(ctx, "rt:"+string(hRefreshToken))
	val := newRefreshTE.Serialize()
	_ = uc.cache.Set(ctx, "at:"+string(newHAccessT), metadataDTO, uc.sDuration)
	_ = uc.cache.Set(ctx, "rt:"+string(newHRefreshT), string(val), uc.rtDuration)

	return &dto.TokenPairDTO{
		RefreshToken: newRefreshT,
		AccessToken:  newAccessT,
	}, nil
}

func (uc *SessionUseCase) RevokeSession(ctx context.Context, sessionID string) error {
	txRepo, tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op if already committed

	// DB Hits
	session, err := txRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	refreshToken, err := txRepo.GetRefreshTokenBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	if err = txRepo.RevokeRefreshToken(ctx, refreshToken.ID()); err != nil {
		return err
	}
	if err = txRepo.DeleteSession(ctx, sessionID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	// removing the cache
	_ = uc.cache.Delete(ctx, "at:"+string(session.Token()))
	_ = uc.cache.Delete(ctx, "rt:"+string(refreshToken.HashedToken()))

	return nil
}

func (uc *SessionUseCase) RevokeAllUserSessions(ctx context.Context, userID string) error {

	sessions, err := uc.repo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var keys []string
	sessionIDs := make([]string, len(sessions))

	for i, s := range sessions {
		sessionIDs[i] = s.ID()
		keys = append(keys, "at:"+string(s.Token()))
	}

	refreshTokens, err := uc.repo.GetActiveRefreshTokensBySessionIDs(ctx, sessionIDs...)
	if err != nil {
		return err
	}
	for _, rt := range refreshTokens {
		keys = append(keys, "rt:"+string(rt.HashedToken()))
	}

	if err = uc.repo.DeleteAllUserSessionsAndRefreshTokens(ctx, userID); err != nil {
		return err
	}
	_ = uc.cache.DeleteMultiple(ctx, keys...)

	return nil
}

func (uc *SessionUseCase) GetUserSessions(ctx context.Context, userID string) ([]*dto.SessionDTO, error) {

	sessions, err := uc.repo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var out []*dto.SessionDTO
	for _, session := range sessions {
		out = append(out, &dto.SessionDTO{
			SessionId: session.ID(),
			Metadata: dto.SessionMetadataDTO{
				IP:          session.IP(),
				ClientAgent: session.UserAgent(),
			},
		})
	}

	return out, nil
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
