package service

import (
	"api/internal/model/domain"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ObjectStorage interface {
	GetPresignedURL(ctx context.Context, pdfKey string, expires time.Duration) (string, error)
	CheckObjectExists(ctx context.Context, pdfKey string) error
}

type JobStore interface {
	CreateJob(ctx context.Context, job domain.Job) error
	MarkQueuedJob(ctx context.Context, job domain.Job) error
	CheckJobStatus(ctx context.Context, jobID uuid.UUID, status string) error
	GetPdfKey(ctx context.Context, jobID uuid.UUID) (string, error)
}

type QueuePublisher interface {
	PublishJob(ctx context.Context, jb domain.Job) error
}

type JobService struct {
	jobStore      JobStore
	objectStorage ObjectStorage
	publisher     QueuePublisher
	logger        *slog.Logger
}

func NewJobService(jobStore JobStore, objectStorage ObjectStorage, publisher QueuePublisher, log *slog.Logger) *JobService {
	return &JobService{
		jobStore:      jobStore,
		objectStorage: objectStorage,
		publisher:     publisher,
		logger:        log,
	}
}

func (s *JobService) InitUpload(ctx context.Context) (*domain.Job, error) {
	const op = "service.JobService.InitUpload"

	jobID := uuid.New()
	pdfKey := fmt.Sprintf("pdf/%s.pdf", jobID)
	now := time.Now()

	uploadURL, err := s.objectStorage.GetPresignedURL(ctx, pdfKey, 5*time.Minute)
	if err != nil {
		s.logger.Debug("failed to get presigned URL", "op", op, "job_id", jobID, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	jb := domain.Job{
		JobID:     jobID,
		Status:    domain.JobStatusCreated,
		PDFKey:    pdfKey,
		UploadURL: uploadURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err = s.jobStore.CreateJob(ctx, jb); err != nil {
		s.logger.Debug("failed to create job", "op", op, "job_id", jobID, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Debug("upload initialised", "op", op, "job_id", jobID)
	return &jb, nil
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID uuid.UUID) error {
	const op = "service.JobService.CompleteUpload"

	if err := s.jobStore.CheckJobStatus(ctx, jobID, "status"); err != nil {
		s.logger.Debug("job status check failed", "op", op, "job_id", jobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	pdfKey, err := s.jobStore.GetPdfKey(ctx, jobID)
	if err != nil {
		s.logger.Debug("failed to get pdf key", "op", op, "job_id", jobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = s.objectStorage.CheckObjectExists(ctx, pdfKey); err != nil {
		s.logger.Debug("object not found in storage", "op", op, "job_id", jobID, "pdf_key", pdfKey, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()

	jb := domain.Job{
		JobID:     jobID,
		Status:    domain.JobStatusQueued,
		PDFKey:    pdfKey,
		UpdatedAt: now,
	}

	if err = s.jobStore.MarkQueuedJob(ctx, jb); err != nil {
		s.logger.Debug("failed to mark job as queued", "op", op, "job_id", jobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = s.publisher.PublishJob(ctx, jb); err != nil {
		s.logger.Debug("failed to publish job", "op", op, "job_id", jobID, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Debug("upload completed", "op", op, "job_id", jobID)
	return nil
}
