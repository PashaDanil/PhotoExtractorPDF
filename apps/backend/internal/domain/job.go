package domain

import (
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

func (j *Job) ToTask() JobTask {
	return JobTask{
		JobID:  j.JobID.String(),
		PDFKey: j.PDFKey,
	}
}

func (j *Job) ToInitResponse() InitResponse {
	return InitResponse{
		JobID:     j.JobID.String(),
		UploadURL: j.UploadURL,
	}
}

func (j *Job) ToRecord() JobRecord {
	return JobRecord{
		JobID:     j.JobID.String(),
		Status:    string(j.Status),
		PDFKey:    j.PDFKey,
		CreatedAt: j.CreatedAt.Unix(),
		UpdatedAt: j.UpdatedAt.Unix(),
	}
}

type JobRecord struct {
	JobID     string
	Status    string
	PDFKey    string
	CreatedAt int64
	UpdatedAt int64
}

type InitResponse struct {
	JobID     string `json:"job_id"`
	UploadURL string `json:"upload_url"`
}

type JobTask struct {
	JobID  string `json:"job_id"`
	PDFKey string `json:"pdf_key"`
}
