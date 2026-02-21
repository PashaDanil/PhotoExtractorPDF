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
	endpoint := cfg.MinIOConfig.URL
	useSSL := cfg.MinIOConfig.UseSSL
	bucket := cfg.MinIOConfig.Bucket

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOConfig.User, cfg.MinIOConfig.Password, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	healthCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	exists, err := client.BucketExists(healthCtx, bucket)
	if err != nil {
		return nil, fmt.Errorf("healthcheck (bucket exists): %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket does not exist: %s", bucket)
	}

	return client, nil
}
