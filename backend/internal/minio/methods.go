package minio

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
}

func NewMinioStorage(client *minio.Client) *MinioStorage {
	return &MinioStorage{
		client: client,
	}
}

func (m *MinioStorage) PresignPut(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, bucketName, objectKey, expires)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *MinioStorage) PresignGet(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectKey, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
