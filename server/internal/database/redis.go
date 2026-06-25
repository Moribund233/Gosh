package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gosh/internal/config"
)

var RedisClient *redis.Client

func InitRedis(cfg config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
