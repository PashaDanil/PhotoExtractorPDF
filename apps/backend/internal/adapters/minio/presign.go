package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageRepo struct {
	client *minio.Client
}

func NewObjectStorageRepo(client *minio.Client) *ObjectStorageRepo {
	return &ObjectStorageRepo{client: client}
}

func (m *ObjectStorageRepo) GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, "imgpdf", pdfKey, expires)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
