package service

import (
	"context"
	"fmt"
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

func (s *JobService) InitUpload(ctx context.Context) (
	jobID string,
	pdfKey string,
	uploadURL string,
	err error,
) {
	jobID = uuid.NewString()
	pdfKey = fmt.Sprintf("pdf/%s.pdf", jobID)

	// Сначала проверяем MinIO - генерируем presigned URL
	uploadURL, err = s.storage.PresignPut(ctx, pdfKey, 15*time.Minute)
	if err != nil {
		return
	}

	// Только если MinIO ОК - создаем запись в Redis
	now := time.Now()
	if err = s.redisRepo.CreateJob(ctx, jobID, pdfKey, now); err != nil {
		return
	}

	return
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID string) error {
	now := time.Now()
	return s.redisRepo.MarkQueuedJob(ctx, jobID, now)
}
