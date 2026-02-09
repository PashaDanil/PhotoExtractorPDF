package minio

import (
	"api/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.MinIO.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.User, cfg.MinIO.Password, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
