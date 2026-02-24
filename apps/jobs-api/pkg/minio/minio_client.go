package minio

import (
	"api/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(ctx context.Context, cfg *config.Config) (*minio.Client, error) {
	const op = "minio.New"

	endpoint := cfg.MinIOConfig.URL
	useSSL := cfg.MinIOConfig.UseSSL
	bucket := cfg.MinIOConfig.Bucket

	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.MinIOConfig.User,
			cfg.MinIOConfig.Password,
			"",
		),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	healthCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	exists, err := client.BucketExists(healthCtx, bucket)
	if err != nil {
		return nil, fmt.Errorf("%s: halthcheck (bucket exists): %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: bucket does not exist: %s", op, bucket)
	}

	return client, nil
}
