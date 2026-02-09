package job

import (
	"context"
	"time"
)

type ObjectStorage interface {
	GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error)
}

type JobStore interface {
	CreateJob(ctx context.Context, job *Job) error
	MarkQueuedJob(ctx context.Context, job *Job) error
}
