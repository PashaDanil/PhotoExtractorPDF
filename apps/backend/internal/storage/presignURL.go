package storage

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type Storage struct {
	client *minio.Client
}

func NewStorage(client *minio.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (m *Storage) PresignPut(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, bucketName, objectKey, expires)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *Storage) PresignGet(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectKey, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
