package redis

import (
	"context"
	"go-api/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type JobRedis struct {
	rdb *redis.Client
}

func NewJobRedis(rdb *redis.Client) *JobRedis {
	return &JobRedis{
		rdb: rdb,
	}
}

func (r *JobRedis) CreateJob(ctx context.Context, job *model.Job) error {
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

func (r *JobRedis) MarkQueuedJob(ctx context.Context, jobID string, now time.Time) error {
	key := "job:" + jobID
	timestamp := now.Unix()

	return r.rdb.HSet(ctx, key, map[string]interface{}{
		"status":     "queued",
		"updated_at": timestamp,
	}).Err()
}
