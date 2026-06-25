package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrNotStored = errors.New("cache set failed")
)

type Cache struct {
	client *redis.Client
	logger *zap.Logger
}

func New(client *redis.Client, logger *zap.Logger) *Cache {
	return &Cache{client: client, logger: logger}
}

func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrCacheMiss) && c.logger != nil {
		c.logger.Warn("cache get failed, falling back", zap.String("key", key), zap.Error(err))
	}

	result, err := fn()
	if err != nil {
		return err
	}

	if err := c.Set(ctx, key, result, ttl); err != nil && c.logger != nil {
		c.logger.Warn("cache set failed", zap.String("key", key), zap.Error(err))
	}

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Cache) Invalidate(ctx context.Context, patterns ...string) error {
	for _, pattern := range patterns {
		iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}
	return nil
}

var defaultCache *Cache

func InitDefault(client *redis.Client, logger *zap.Logger) {
	defaultCache = New(client, logger)
}

func Default() *Cache {
	return defaultCache
}
