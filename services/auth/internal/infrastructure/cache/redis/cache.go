package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aligh5331/godrop/services/auth/internal/domain"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// Set serializes value to JSON and stores it with an optional TTL.
// Pass 0 for expiration to store without TTL.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get returns the raw JSON string stored at key.
// Returns an error wrapping ErrNotFound if the key doesn't exist.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrCacheNotFound
	}
	return val, err
}

// Exists returns true if the key is present in Redis.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Delete removes a single key. Does not error if the key doesn't exist.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeleteMultiple removes multiple keys in a single round-trip.
func (c *Cache) DeleteMultiple(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
