package repository

import (
	"context"
	"time"
)

type JobRedisRepo interface {
	CreateJob(ctx context.Context, jobID string, pdfKey string, now time.Time) error
	MarkQueuedJob(ctx context.Context, jobID string, now time.Time) error
}
