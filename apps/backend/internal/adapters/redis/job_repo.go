package redis

import (
	"api/internal/domain/job"
	"api/pkg/errorx"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type JobStoreRepo struct {
	rdb *redis.Client
}

func NewJobStoreRepo(rdb *redis.Client) *JobStoreRepo {
	return &JobStoreRepo{rdb: rdb}
}

func (r *JobStoreRepo) CreateJob(ctx context.Context, jb *job.Job) error {
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

func (r *JobStoreRepo) MarkQueuedJob(ctx context.Context, jb *job.Job) error {
	currentStatus, err := r.rdb.HGet(ctx, jb.JobID, "status").Result()
	if err != nil {
		// если пусто по ключу ErrNotFound 404
		if err == redis.Nil {
			return fmt.Errorf("job %s not found: %w", jb.JobID, errorx.ErrNotFound)
		}
		return err
	}
	// если уже queued ErrAlreadyCompleted 409
	if currentStatus == string(job.JobStatusQueued) {
		return fmt.Errorf("job %s already queued: %w", jb.JobID, errorx.ErrAlreadyCompleted)
	}

	cmd := r.rdb.HSet(ctx, jb.JobID, map[string]any{
		"status":     string(jb.Status),
		"updated_at": jb.UpdatedAt,
	})
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
