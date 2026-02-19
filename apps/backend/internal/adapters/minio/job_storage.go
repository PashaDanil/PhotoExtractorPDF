package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageRepo struct {
	client *minio.Client
	bucket string
}

func NewObjectStorageRepo(mio *minio.Client, bucket string) *ObjectStorageRepo {
	return &ObjectStorageRepo{client: mio, bucket: bucket}
}

func (m *ObjectStorageRepo) GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error) {
	const op = "ObjectStorageRepo.GetPresignedURL"

	u, err := m.client.PresignedPutObject(ctx, m.bucket, pdfKey, expires)
	if err != nil {
		// обработать ошибку
		return "", err
	}

	return u.String(), nil
}

func (m *ObjectStorageRepo) CheckObjectExists(ctx context.Context, pdfKey string) error {
	const op = "ObjectStorageRepo.CheckObjectExists"

	_, err := m.client.StatObject(ctx, m.bucket, pdfKey, minio.StatObjectOptions{})
	if err != nil {
		// обработать ошибку
		return err
	}

	return nil
}
