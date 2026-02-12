package job

import (
	"context"
	"time"
)

type ObjectStorage interface {
	GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error)
	CheckObjectExists(ctx context.Context, pdfKey string) error
}

type JobStore interface {
	CreateJob(ctx context.Context, job *Job) error
	MarkQueuedJob(ctx context.Context, job *Job) error
	CheckJobStatusQueued(ctx context.Context, jobID string) error
	GetPdfKey(ctx context.Context, jobID string) (string, error)
}
