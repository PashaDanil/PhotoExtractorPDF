package redis

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

func normalizeRedisErr(op string, err error) error {
	if err == nil {
		return nil
	}

	// Контекст — всегда наверх как есть
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// redis.Nil обычно обрабатывается в конкретном методе (NotFound)
	if errors.Is(err, redis.Nil) {
		return err
	}

	// Сетевые/таймауты
	var ne net.Error
	if errors.As(err, &ne) {
		if ne.Timeout() {
			return fmt.Errorf("%s: redis timeout: %w", op, err)
		}
		return fmt.Errorf("%s: redis network error: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, err)
}
