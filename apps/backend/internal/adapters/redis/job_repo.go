package redis

import (
	"api/internal/domain/job"
	"context"

	"github.com/redis/go-redis/v9"
)

type JobStoreRepo struct {
	rdb *redis.Client
}

func NewJobStoreRepo(rdb *redis.Client) *JobStoreRepo {
	return &JobStoreRepo{rdb: rdb}
}

func (r *JobStoreRepo) CreateJob(ctx context.Context, job *job.Job) error {
	key := "job:" + job.JobID

	err := r.rdb.HSet(ctx, key, map[string]any{
		"status":     job.Status,
		"pdf_key":    job.PDFKey,
		"upload_url": job.UploadURL,
		"created_at": job.CreatedAt,
		"updated_at": job.UpdatedAt,
	})
	if err != nil {
		return err.Err()
	}

	return nil
}

func (r *JobStoreRepo) MarkQueuedJob(ctx context.Context, job *job.Job) error {
	key := "job:" + job.JobID

	err := r.rdb.HSet(ctx, key, map[string]any{
		"status":     job.Status,
		"updated_at": job.UpdatedAt,
	})
	if err != nil {
		return err.Err()
	}

	return nil
}
