package redis

import (
	"context"
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

func (r *JobRedis) CreateJob(ctx context.Context, jobID string, pdfKey string, now time.Time) error {
	key := "job:" + jobID
	timestamp := now.Unix()

	return r.rdb.HSet(ctx, key, map[string]interface{}{
		"status":     "created",
		"pdf_key":    pdfKey,
		"created_at": timestamp,
		"updated_at": timestamp,
	}).Err()
}

func (r *JobRedis) MarkQueuedJob(ctx context.Context, jobID string, now time.Time) error {
	key := "job:" + jobID
	timestamp := now.Unix()

	return r.rdb.HSet(ctx, key, map[string]interface{}{
		"status":     "queued",
		"updated_at": timestamp,
	}).Err()
}
