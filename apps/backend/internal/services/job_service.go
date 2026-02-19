package services

import (
	"api/internal/domain/job"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	jobStore      job.JobStore
	objectStorage job.ObjectStorage
	publisher     job.QueuePublisher
}

func NewJobService(jobStore job.JobStore, objectStorage job.ObjectStorage, publisher job.QueuePublisher) *JobService {
	return &JobService{
		jobStore:      jobStore,
		objectStorage: objectStorage,
		publisher:     publisher,
	}
}

func (s *JobService) InitUpload(ctx context.Context) (*job.Job, string, error) {
	jobID := uuid.NewString()
	pdfKey := fmt.Sprintf("pdf/%s.pdf", jobID)
	now := time.Now()

	uploadURL, err := s.objectStorage.GetPresignedURL(ctx, pdfKey, 5*time.Minute)
	if err != nil {
		// обработать ошибку
		return nil, "", err
	}

	jb := &job.Job{
		JobID:     jobID,
		Status:    job.JobStatusCreated,
		PDFKey:    pdfKey,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	if err = s.jobStore.CreateJob(ctx, jb); err != nil {
		// обработать ошибку
		return nil, "", err
	}

	return jb, uploadURL, nil
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID string) error {
	// проверяем статус
	err := s.jobStore.CheckJobStatusQueued(ctx, jobID)
	if err != nil {
		// обработать ошибку
		return err
	}

	// получаем pdfKey
	pdfKey, err := s.jobStore.GetPdfKey(ctx, jobID)
	if err != nil {
		// обработать ошибку
		return err
	}

	// проверяем, что объект существует в s3 хранилище
	err = s.objectStorage.CheckObjectExists(ctx, pdfKey)
	if err != nil {
		// обработать ошибку
		return err
	}

	now := time.Now()

	jb := &job.Job{
		JobID:     jobID,
		Status:    job.JobStatusQueued,
		UpdatedAt: now.Unix(),
	}

	err = s.jobStore.MarkQueuedJob(ctx, jb)
	if err != nil {
		// обработать ошибку
		return err
	}

	// публикуем задачу в очередь
	err = s.publisher.PublishJob(ctx, jobID, pdfKey)
	if err != nil {
		// обработать ошибку
		return err
	}

	return nil
}
