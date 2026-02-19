package minio

import (
	"api/internal/domain/job"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/minio/minio-go/v7"
)

func normalizeMinioErr(op string, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	var ne net.Error
	if errors.As(err, &ne) {
		if ne.Timeout() {
			return fmt.Errorf("%s: storage timeout: %w", op, err) // замапить в 503
		}
		return fmt.Errorf("%s: storage network error: %w", op, err)
	}

	resp := minio.ToErrorResponse(err)
	if resp.Code != "" || resp.StatusCode != 0 {

		if resp.Code == "NoSuchKey" || resp.Code == "NoSuchObject" || resp.StatusCode == 404 {
			return job.ErrObjectNotFound
		}

		if resp.Code == "AccessDenied" || resp.StatusCode == 403 {
			return job.ErrStorageForbidden
		}
		if resp.StatusCode == 429 || resp.StatusCode == 503 {
			return job.ErrStorageUnavailable
		}

		return fmt.Errorf("%s: storage error (code=%s status=%d): %w", op, resp.Code, resp.StatusCode, err)
	}

	return fmt.Errorf("%s: %w", op, err)
}
