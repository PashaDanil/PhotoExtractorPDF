package model

type JobStatus string

const (
	JobStatusCreated    JobStatus = "created"
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusDone       JobStatus = "done"
)

type Job struct {
	JobID     string
	Status    JobStatus
	PDFKey    string
	UploadURL string
	CreatedAt int64
	UpdatedAt int64
}
