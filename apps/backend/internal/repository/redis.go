package repository

import (
	"context"
	"go-api/internal/model"
	"time"
)

type JobRedisRepo interface {
	CreateJob(ctx context.Context, job *model.Job) error
	MarkQueuedJob(ctx context.Context, jobID string, now time.Time) error
}
