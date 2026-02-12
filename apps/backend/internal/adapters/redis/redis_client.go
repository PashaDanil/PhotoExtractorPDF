package redis

import (
	"api/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg *config.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Redis{client: rdb}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	return r.client.Close()
}
