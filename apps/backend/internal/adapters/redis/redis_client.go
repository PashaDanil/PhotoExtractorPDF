package redis

import (
	"api/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(
	ctx context.Context,
	cfg *config.Config,
) (*Redis, error) {
	addr := cfg.RedisConfig.URL
	db := cfg.RedisConfig.DB

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisConfig.Password,
		DB:       db,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("%s: ping failed: %w", err)
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	if r.client == nil {
		return nil
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("%s: close client: %w", err)
	}

	return nil
}
