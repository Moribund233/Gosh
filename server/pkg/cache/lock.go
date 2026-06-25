package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const lockPrefix = "lock:"

type Lock struct {
	client *redis.Client
}

func NewLock(client *redis.Client) *Lock {
	return &Lock{client: client}
}

func (l *Lock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := l.client.SetNX(ctx, lockPrefix+key, 1, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (l *Lock) Release(ctx context.Context, key string) error {
	return l.client.Del(ctx, lockPrefix+key).Err()
}
