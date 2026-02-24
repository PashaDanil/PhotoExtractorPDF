package cache

import (
	errs "api/internal/errors"
	"api/internal/model/domain"
	redisClient "api/pkg/redis"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JobStoreRepo struct {
	rdb    *redis.Client
	logger *slog.Logger
}

func NewJobStoreRepo(r *redisClient.Redis, log *slog.Logger) *JobStoreRepo {
	return &JobStoreRepo{rdb: r.Client(), logger: log}
}

func (r *JobStoreRepo) CreateJob(ctx context.Context, jb domain.Job) error {
	const op = "cache.JobStoreRepo.CreateJob"

	record := jb.ToRecord()

	r.logger.Debug("creating job in cache", "op", op, "job_id", record.JobID)

	cmd := r.rdb.HSet(ctx, record.JobID, map[string]any{
		"status":     record.Status,
		"pdf_key":    record.PDFKey,
		"created_at": record.CreatedAt,
		"updated_at": record.UpdatedAt,
	})
	if err := cmd.Err(); err != nil {
		r.logger.Debug("failed to create job in cache", "op", op, "job_id", record.JobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	r.logger.Debug("job created in cache", "op", op, "job_id", record.JobID)
	return nil
}

func (r *JobStoreRepo) MarkQueuedJob(ctx context.Context, jb domain.Job) error {
	const op = "cache.JobStoreRepo.MarkQueuedJob"

	record := jb.ToRecord()

	r.logger.Debug("marking job as queued", "op", op, "job_id", record.JobID)

	cmd := r.rdb.HSet(ctx, record.JobID, map[string]any{
		"status":     record.Status,
		"updated_at": record.UpdatedAt,
	})
	if err := cmd.Err(); err != nil {
		r.logger.Debug("failed to mark job as queued", "op", op, "job_id", record.JobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	r.logger.Debug("job marked as queued", "op", op, "job_id", record.JobID)
	return nil
}

func (r *JobStoreRepo) CheckJobStatus(ctx context.Context, jobID uuid.UUID, status string) error {
	const op = "cache.JobStoreRepo.CheckJobStatus"

	key := jobID.String()

	r.logger.Debug("checking job status", "op", op, "job_id", key)

	currentStatus, err := r.rdb.HGet(ctx, key, "status").Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("job not found in cache", "op", op, "job_id", key)
			return errs.ErrNotFound
		}
		r.logger.Debug("failed to get job status", "op", op, "job_id", key, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	r.logger.Debug("got job status", "op", op, "job_id", key, "status", currentStatus)

	if currentStatus == string(domain.JobStatusQueued) {
		return errs.ErrAlreadyQueued
	}
	return nil
}

func (r *JobStoreRepo) GetPdfKey(ctx context.Context, jobID uuid.UUID) (string, error) {
	const op = "cache.JobStoreRepo.GetPdfKey"

	key := jobID.String()

	r.logger.Debug("getting pdf key", "op", op, "job_id", key)

	pdfKey, err := r.rdb.HGet(ctx, key, "pdf_key").Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("job not found in cache", "op", op, "job_id", key)
			return "", errs.ErrNotFound
		}
		r.logger.Debug("failed to get pdf key", "op", op, "job_id", key, "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	r.logger.Debug("got pdf key", "op", op, "job_id", key)
	return pdfKey, nil
}
