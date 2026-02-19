package minio

import (
	"api/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) (*minio.Client, error) {
	const op = "minio.New"

	endpoint := cfg.MinIOConfig.URL
	useSSL := cfg.MinIOConfig.UseSSL
	bucket := cfg.MinIOConfig.Bucket

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOConfig.User, cfg.MinIOConfig.Password, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: create client: %w", op, err)
	}

	healthCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	exists, err := client.BucketExists(healthCtx, bucket)
	if err != nil {
		return nil, fmt.Errorf("%s: healthcheck (bucket exists): %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: bucket does not exist: %s", op, bucket)
	}

	log.Info("minio ready",
		slog.String("component", "minio"),
		slog.String("endpoint", endpoint),
		slog.Bool("ssl", useSSL),
		slog.String("bucket", bucket),
	)

	return client, nil
}
