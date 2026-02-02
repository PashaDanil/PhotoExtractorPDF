package service

import (
	"context"
	"fmt"
	"imgpdf/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	repoRedis    repository.JobRedisRepo
	storageMinIO repository.MinioRepo
}

func NewJobService(repoRedis repository.JobRedisRepo, storageMinIO repository.MinioRepo) *JobService {
	return &JobService{
		repoRedis:    repoRedis,
		storageMinIO: storageMinIO,
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
	uploadURL, err = s.storageMinIO.PresignPut(ctx, pdfKey, 15*time.Minute)
	if err != nil {
		return
	}

	// Только если MinIO ОК - создаем запись в Redis
	now := time.Now()
	if err = s.repoRedis.CreateJob(ctx, jobID, pdfKey, now); err != nil {
		return
	}

	return
}

func (s *JobService) CompleteUpload(ctx context.Context, jobID string) error {
	now := time.Now()
	return s.repoRedis.MarkQueuedJob(ctx, jobID, now)
}
