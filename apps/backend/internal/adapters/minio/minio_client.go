package minio

import (
	"api/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.MinIOConfig.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOConfig.User, cfg.MinIOConfig.Password, ""),
		Secure: cfg.MinIOConfig.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
