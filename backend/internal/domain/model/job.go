package model

type JobStatus string

const (
	JobStatusCreated    JobStatus = "created"
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusDone       JobStatus = "done"
	JobStatusError      JobStatus = "error"
)

type Job struct {
	Status    JobStatus `json:"status"`
	PDFKey    string    `json:"pdf_key"`
	ZIPKey    string    `json:"zip_key,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}
