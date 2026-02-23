package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/domain"
	"github.com/aligh5331/godrop/services/auth/internal/domain/entity"
	"github.com/aligh5331/godrop/services/auth/internal/domain/helpers"
	"github.com/aligh5331/godrop/services/auth/internal/domain/repository"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repository.SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

type gSession struct {
	ID          string `gorm:"primaryKey"`
	UserID      string
	UserAgent   string
	IP          string `gorm:"column:ip"`
	HashedToken string
	IsRevoked   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}

func (*gSession) TableName() string {
	return "sessions"
}
func toGSession(s *entity.Session) *gSession {
	return &gSession{
		ID:          s.ID(),
		UserID:      s.UserID(),
		UserAgent:   s.UserAgent(),
		IP:          s.IP(),
		HashedToken: string(s.Token()),
		IsRevoked:   s.IsRevoke(),
		CreatedAt:   s.CreatedAt(),
		UpdatedAt:   s.UpdatedAt(),
		ExpiresAt:   s.ExpiresAt(),
	}
}

func (gs *gSession) toDSession() (*entity.Session, error) {
	return entity.NewSession(
		gs.ID,
		gs.UserID,
		gs.UserAgent,
		gs.IP,
		entity.HashedToken(gs.HashedToken),
		gs.CreatedAt,
		gs.UpdatedAt,
		gs.ExpiresAt,
	)
}

type gRefreshToken struct {
	ID          string `gorm:"primaryKey"`
	SessionID   string
	HashedToken string
	CreatedAt   time.Time
	RevokedAt   time.Time
	ExpiresAt   time.Time
}

func (*gRefreshToken) TableName() string {
	return "refresh_tokens"
}
func toGRefreshToken(r *entity.RefreshToken) *gRefreshToken {
	return &gRefreshToken{
		ID:          r.ID(),
		SessionID:   r.SessionID(),
		HashedToken: string(r.HashedToken()),
		CreatedAt:   r.CreatedAt(),
		RevokedAt:   r.RevokedAt(),
		ExpiresAt:   r.ExpiresAt(),
	}
}
func (gr *gRefreshToken) toDRefreshToken() (*entity.RefreshToken, error) {
	return entity.NewRefreshToken(
		gr.ID,
		gr.SessionID,
		entity.HashedToken(gr.HashedToken),
		gr.CreatedAt,
		gr.ExpiresAt,
		gr.RevokedAt,
	)
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	db := GetDB(ctx, r.db)
	return db.WithContext(ctx).Create(toGSession(session)).Error
}

func (r *SessionRepository) CreateRefreshToken(ctx context.Context, refreshToken *entity.RefreshToken) error {
	db := GetDB(ctx, r.db)

	return db.WithContext(ctx).Create(toGRefreshToken(refreshToken)).Error
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	var session *gSession
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return session.toDSession()
}

func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, token entity.HashedToken) (*entity.Session, error) {
	var refreshToken *gRefreshToken
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Where("hashed_token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return r.GetSessionByID(ctx, refreshToken.SessionID)
}

func (r *SessionRepository) GetSessionByAccessToken(ctx context.Context, token entity.HashedToken) (*entity.Session, error) {
	var session *gSession
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Where("hashed_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return session.toDSession()
}

func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	var models []gSession
	db := GetDB(ctx, r.db)

	result := db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ? AND is_revoked = false", userID, time.Now().UTC()).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	out := make([]*entity.Session, len(models))
	for i, m := range models {
		out[i], _ = m.toDSession()
	}
	return out, nil
}
func (r *SessionRepository) GetRefreshTokenBySessionID(ctx context.Context, sessionID string) (*entity.RefreshToken, error) {
	var refreshToken *gRefreshToken
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Where("session_id = ? AND revoked_at = ?", sessionID, time.Time{}).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return refreshToken.toDRefreshToken()
}

func (r *SessionRepository) GetRefreshTokenEntityByRefreshToken(ctx context.Context, token entity.HashedToken) (*entity.RefreshToken, error) {
	var refreshToken *gRefreshToken
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Where("hashed_token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return refreshToken.toDRefreshToken()
}

func (r *SessionRepository) GetActiveRefreshTokensBySessionIDs(ctx context.Context, sessionIDs ...string) ([]*entity.RefreshToken, error) {
	var models []gRefreshToken
	db := GetDB(ctx, r.db)

	result := db.WithContext(ctx).
		Where("session_id IN ? AND revoked_at = ?", sessionIDs, time.Time{}).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	out := make([]*entity.RefreshToken, len(models))
	for i, m := range models {
		out[i], _ = m.toDRefreshToken()
	}
	return out, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *entity.Session) error {
	db := GetDB(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&gSession{}).
		Where("id = ?", session.ID()).
		Updates(map[string]interface{}{
			"user_agent":   session.UserAgent(),
			"ip":           session.IP(),
			"hashed_token": session.Token(),
			"is_revoked":   session.IsRevoke(),
			"updated_at":   session.UpdatedAt(),
			"expires_at":   session.ExpiresAt(),
		})

	if res.Error != nil {
		return fmt.Errorf("gorm update session: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *SessionRepository) RevokeRefreshToken(ctx context.Context, refreshTokenID string) error {
	db := GetDB(ctx, r.db)

	res := db.WithContext(ctx).
		Model(&gRefreshToken{}).
		Where("id = ?", refreshTokenID).
		Updates(map[string]interface{}{
			"revoked_at": time.Now().UTC(),
		})
	if res.RowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}
	return res.Error
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	db := GetDB(ctx, r.db)

	res := db.WithContext(ctx).
		Where("id = ?", sessionID).
		Delete(&gSession{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *SessionRepository) DeleteAllUserSessionsAndRefreshTokens(ctx context.Context, userID string) error {
	db := GetDB(ctx, r.db)

	// delete RTs via subquery
	err := db.WithContext(ctx).
		Where("session_id IN (?)",
			db.Model(&gSession{}).
				Select("id").
				Where("user_id = ?", userID),
		).
		Delete(&gRefreshToken{}).Error
	if err != nil {
		return err
	}

	return db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&gSession{}).Error
}

type SRTransaction struct {
	tx *gorm.DB
}

func (t *SRTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *SRTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (r *SessionRepository) BeginTx(ctx context.Context) (context.Context, repository.Transaction, error) {
	if _, ok := ctx.Value(helpers.TxKey{}).(*gorm.DB); ok {
		// A transaction already exists. Return the context as-is
		// and a NO-OP transaction so the caller doesn't break the parent TX.
		return ctx, &noOpTx{}, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	newCtx := context.WithValue(ctx, helpers.TxKey{}, tx)
	return newCtx, &SRTransaction{tx: tx}, nil
}
