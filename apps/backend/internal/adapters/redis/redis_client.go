package redis

import (
	"api/pkg/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg *config.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.URL,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})

	return &Redis{client: rdb}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	return r.client.Close()
}
