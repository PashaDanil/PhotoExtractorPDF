package redis

import (
	"api/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) (*Redis, error) {
	const op = "redis.New"

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
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	log.Info("redis ready",
		slog.String("component", "redis"),
		slog.String("addr", addr),
		slog.Int("db", db),
	)

	return &Redis{client: rdb}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	const op = "redis.Close"

	if r.client == nil {
		return nil
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("%s: close client: %w", op, err)
	}

	return nil
}
