package domain

import (
	"api/internal/model/dto"
	"api/internal/model/records"
	"api/internal/model/tasks"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	JobID     uuid.UUID
	Status    JobStatus
	PDFKey    string
	UploadURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobStatus string

const (
	JobStatusCreated    JobStatus = "uploading"
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusProcessed  JobStatus = "processed"
	JobStatusDone       JobStatus = "done"
)

func (j *Job) ToTask() tasks.JobTask {
	return tasks.JobTask{
		JobID:  j.JobID.String(),
		PDFKey: j.PDFKey,
	}
}

func (j *Job) ToInitResponse() dto.InitResponse {
	return dto.InitResponse{
		JobID:     j.JobID.String(),
		UploadURL: j.UploadURL,
	}
}

func (j *Job) ToRecord() records.JobRecord {
	return records.JobRecord{
		JobID:     j.JobID.String(),
		Status:    string(j.Status),
		PDFKey:    j.PDFKey,
		CreatedAt: j.CreatedAt.Unix(),
		UpdatedAt: j.UpdatedAt.Unix(),
	}
}
