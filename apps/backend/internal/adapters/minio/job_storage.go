package minio

import (
	"api/pkg/errorx"
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageRepo struct {
	client *minio.Client
}

func NewObjectStorageRepo(mio *minio.Client) *ObjectStorageRepo {
	return &ObjectStorageRepo{client: mio}
}

func (m *ObjectStorageRepo) GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, "imgpdf", pdfKey, expires)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *ObjectStorageRepo) CheckObjectExists(ctx context.Context, pdfKey string) error {
	info, err := m.client.StatObject(ctx, "imgpdf", pdfKey, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return errorx.ErrObjectNotFound
		}
		return err
	}

	if info.Size <= 0 {
		return errorx.ErrObjectNotFound
	}

	return nil
}
