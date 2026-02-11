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
}

func NewJobService(jobStore job.JobStore, objectStorage job.ObjectStorage) *JobService {
	return &JobService{
		jobStore:      jobStore,
		objectStorage: objectStorage,
	}
}

func (s *JobService) InitUpload(ctx context.Context) (*job.Job, error) {
	jobID := uuid.NewString()
	pdfKey := fmt.Sprintf("pdf/%s.pdf", jobID)
	now := time.Now()

	uploadURL, err := s.objectStorage.GetPresignedURL(ctx, pdfKey, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	jb := &job.Job{
		JobID:     jobID,
		Status:    job.JobStatusCreated,
		PDFKey:    pdfKey,
		UploadURL: uploadURL,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	if err = s.jobStore.CreateJob(ctx, jb); err != nil {
		return nil, err
	}

	return jb, nil
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID string) error {
	now := time.Now()

	jb := &job.Job{
		JobID:     jobID,
		Status:    job.JobStatusQueued,
		UpdatedAt: now.Unix(),
	}

	err := s.jobStore.MarkQueuedJob(ctx, jb)
	if err != nil {
		return err
	}

	return nil
}
