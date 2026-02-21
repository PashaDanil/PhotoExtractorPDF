package redis

import (
	"api/internal/domain"
	errs "api/internal/errors"
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JobStoreRepo struct {
	rdb *redis.Client
}

func NewJobStoreRepo(r *Redis) *JobStoreRepo {
	return &JobStoreRepo{rdb: r.Client()}
}

func (r *JobStoreRepo) CreateJob(ctx context.Context, jb *domain.Job) error {
	key := jb.JobID.String()

	err := r.rdb.HSet(ctx, key, map[string]any{
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
	key := jb.JobID.String()

	err := r.rdb.HSet(ctx, key, map[string]any{
		"status":     string(jb.Status),
		"updated_at": jb.UpdatedAt,
	})
	if err != nil {
		return err.Err()
	}

	return nil
}

func (r *JobStoreRepo) CheckJobStatusQueued(ctx context.Context, jobID uuid.UUID) error {
	key := jobID.String()

	currentStatus, err := r.rdb.HGet(ctx, key, "status").Result()
	if err != nil {
		return err
	}
	if currentStatus == string(domain.JobStatusQueued) {
		return errs.ErrAlreadyCompleted
	}
	return nil
}

func (r *JobStoreRepo) GetPdfKey(ctx context.Context, jobID uuid.UUID) (string, error) {
	key := jobID.String()

	pdfKey, err := r.rdb.HGet(ctx, key, "pdf_key").Result()
	if err != nil {
		if err == redis.Nil {
			return "", errs.ErrNotFound
		}
		return "", err
	}
	return pdfKey, nil
}
