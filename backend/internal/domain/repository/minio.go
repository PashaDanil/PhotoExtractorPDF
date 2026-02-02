package repository

import (
	"context"
	"time"
)

type MinioRepo interface {
	// что-то типо такого
	PresignPut(ctx context.Context, objectKey string, expires time.Duration) (string, error)
	PresignGet(ctx context.Context, objectKey string, expires time.Duration) (string, error)
}
