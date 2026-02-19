package minio

import (
	"api/internal/config"
	"context"
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
		Creds:  credentials.NewStaticV4(cfg.MinIOConfig.User, cfg.MinIOConfig.Password, ""),
		Secure: useSSL,
	})
	if err != nil {
		// обработать ошибку
		return nil, err
	}

	healthCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	exists, err := client.BucketExists(healthCtx, bucket)
	if err != nil {
		// обработать ошибку
		return nil, err
	}
	if !exists {
		// обработать ошибку
		return nil, err
	}

	return client, nil
}
