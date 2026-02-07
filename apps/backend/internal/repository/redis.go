package repository

import (
	"context"
	"go-api/internal/model"
)

type JobRedisRepo interface {
	CreateJob(ctx context.Context, job *model.Job) error
	MarkQueuedJob(ctx context.Context, job *model.Job) error
}
