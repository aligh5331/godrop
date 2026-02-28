package gorm

import (
	"context"

	"github.com/aligh5331/godrop/services/core-api/internal/domain/repo"
	"gorm.io/gorm"
)

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(repo.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}

type noOpTx struct{}

func (n *noOpTx) Commit() error   { return nil } // Success! (But does nothing)
func (n *noOpTx) Rollback() error { return nil } // Safe! (But does nothing)
