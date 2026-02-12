package job

type Job struct {
	JobID     string    `json:"job_id"`
	Status    JobStatus `json:"status"`
	PDFKey    string    `json:"pdf_key"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type JobStatus string

const (
	JobStatusCreated    JobStatus = "uploading"
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusDone       JobStatus = "done"
)
