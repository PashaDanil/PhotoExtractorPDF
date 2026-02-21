package redis

import (
	"api/internal/domain"
	errs "api/internal/errors"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type JobStoreRepo struct {
	rdb *redis.Client
}

func NewJobStoreRepo(r *Redis) *JobStoreRepo {
	return &JobStoreRepo{rdb: r.Client()}
}

func (r *JobStoreRepo) CreateJob(ctx context.Context, jb *domain.Job) error {
	err := r.rdb.HSet(ctx, jb.JobID, map[string]any{
		"status":     string(jb.Status),
		"pdf_key":    jb.PDFKey,
		"created_at": jb.CreatedAt,
		"updated_at": jb.UpdatedAt,
	})
	if err != nil {
		return err.Err()
	}

	return nil
}

func (r *JobStoreRepo) MarkQueuedJob(ctx context.Context, jb *domain.Job) error {
	cmd := r.rdb.HSet(ctx, jb.JobID, map[string]any{
		"status":     string(jb.Status),
		"updated_at": jb.UpdatedAt,
	})
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *JobStoreRepo) CheckJobStatusQueued(ctx context.Context, jobID string) error {
	currentStatus, err := r.rdb.HGet(ctx, jobID, "status").Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("job %s not found: %w", jobID, errs.ErrNotFound)
		}
		return err
	}
	if currentStatus == string(domain.JobStatusQueued) {
		return fmt.Errorf("job %s already queued: %w", jobID, errs.ErrAlreadyCompleted)
	}
	return nil
}

func (r *JobStoreRepo) GetPdfKey(ctx context.Context, jobID string) (string, error) {
	pdfKey, err := r.rdb.HGet(ctx, jobID, "pdf_key").Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("job %s not found: %w", jobID, errs.ErrNotFound)
		}
		return "", err
	}
	return pdfKey, nil
}
