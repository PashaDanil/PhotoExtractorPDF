package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Addr: os.Getenv("REDIS_URL"),
		// Password: "redis",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
