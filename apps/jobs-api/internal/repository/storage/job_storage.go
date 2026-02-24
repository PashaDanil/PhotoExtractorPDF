package storage

import (
	errs "api/internal/errors"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStorageRepo struct {
	client *minio.Client
	bucket string
	logger *slog.Logger
}

func NewObjectStorageRepo(mio *minio.Client, bucket string, log *slog.Logger) *ObjectStorageRepo {
	return &ObjectStorageRepo{client: mio, bucket: bucket, logger: log}
}

func (m *ObjectStorageRepo) GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error) {
	const op = "storage.ObjectStorageRepo.GetPresignedURL"

	m.logger.Debug("generating presigned URL", "op", op, "pdf_key", pdfKey, "expires", expires)

	u, err := m.client.PresignedPutObject(ctx, m.bucket, pdfKey, expires)
	if err != nil {
		m.logger.Debug("failed to generate presigned URL", "op", op, "pdf_key", pdfKey, "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	m.logger.Debug("presigned URL generated", "op", op, "pdf_key", pdfKey)
	return u.String(), nil
}

func (m *ObjectStorageRepo) CheckObjectExists(ctx context.Context, pdfKey string) error {
	const op = "storage.ObjectStorageRepo.CheckObjectExists"

	m.logger.Debug("checking object existence", "op", op, "pdf_key", pdfKey)

	_, err := m.client.StatObject(ctx, m.bucket, pdfKey, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			m.logger.Debug("object not found", "op", op, "pdf_key", pdfKey)
			return errs.ErrNotFound
		}
		m.logger.Debug("failed to stat object", "op", op, "pdf_key", pdfKey, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	m.logger.Debug("object exists", "op", op, "pdf_key", pdfKey)
	return nil
}
