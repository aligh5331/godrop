package repository

import (
	"context"
	"time"
)

type CacheRepository interface {
	// Set with TTL (Time-To-Live)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
