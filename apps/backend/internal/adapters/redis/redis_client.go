package redis

import (
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redis",
		DB:       0,
	})

	return rdb, nil
}
