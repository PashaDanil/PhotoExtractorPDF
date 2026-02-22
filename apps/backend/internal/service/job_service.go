package service

import (
	"api/internal/model/domain"
	"context"
	"fmt"
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
	CheckJobStatusQueued(ctx context.Context, jobID uuid.UUID) error
	GetPdfKey(ctx context.Context, jobID uuid.UUID) (string, error)
}

type QueuePublisher interface {
	PublishJob(ctx context.Context, jb domain.Job) error
}

type JobService struct {
	jobStore      JobStore
	objectStorage ObjectStorage
	publisher     QueuePublisher
}

func NewJobService(jobStore JobStore, objectStorage ObjectStorage, publisher QueuePublisher) *JobService {
	return &JobService{
		jobStore:      jobStore,
		objectStorage: objectStorage,
		publisher:     publisher,
	}
}

func (s *JobService) InitUpload(ctx context.Context) (*domain.Job, error) {
	jobID := uuid.New()
	pdfKey := fmt.Sprintf("pdf/%s.pdf", jobID)
	now := time.Now()

	uploadURL, err := s.objectStorage.GetPresignedURL(ctx, pdfKey, 5*time.Minute)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &jb, nil
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID uuid.UUID) error {
	// проверяем статус
	err := s.jobStore.CheckJobStatusQueued(ctx, jobID)
	if err != nil {
		return err
	}

	// получаем pdfKey
	pdfKey, err := s.jobStore.GetPdfKey(ctx, jobID)
	if err != nil {
		return err
	}

	// проверяем, что объект существует в s3 хранилище
	err = s.objectStorage.CheckObjectExists(ctx, pdfKey)
	if err != nil {
		return err
	}

	now := time.Now()

	jb := domain.Job{
		JobID:     jobID,
		Status:    domain.JobStatusQueued,
		PDFKey:    pdfKey,
		UpdatedAt: now,
	}

	err = s.jobStore.MarkQueuedJob(ctx, jb)
	if err != nil {
		return err
	}

	// публикуем задачу в очередь
	err = s.publisher.PublishJob(ctx, jb)
	if err != nil {
		return err
	}

	return nil
}
