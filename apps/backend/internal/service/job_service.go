package service

import (
	"context"
	"fmt"
	"go-api/internal/model"
	"go-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	redisRepo repository.JobRedisRepo
	storage   repository.StorageRepo
}

func NewJobService(redisRepo repository.JobRedisRepo, storage repository.StorageRepo) *JobService {
	return &JobService{
		redisRepo: redisRepo,
		storage:   storage,
	}
}

func (s *JobService) InitUpload(ctx context.Context) (*model.Job, error) {
	jobID := uuid.NewString()
	pdfKey := fmt.Sprintf("pdf/%s.pdf", jobID)
	now := time.Now()

	uploadURL, err := s.storage.PresignPut(ctx, pdfKey, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	job := &model.Job{
		JobID:     jobID,
		Status:    model.JobStatusCreated,
		PDFKey:    pdfKey,
		UploadURL: uploadURL,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	if err = s.redisRepo.CreateJob(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID string) error {
	now := time.Now()
	return s.redisRepo.MarkQueuedJob(ctx, jobID, now)
}
