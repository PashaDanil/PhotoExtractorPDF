package job

type Job struct {
	JobID     string    `json:"job_id"`
	Status    JobStatus `json:"status"`
	PDFKey    string    `json:"pdf_key"`
	UploadURL string    `json:"upload_url"`
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

// MarshalBinary implements encoding.BinaryMarshaler for Redis compatibility
func (s JobStatus) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Redis compatibility
func (s *JobStatus) UnmarshalBinary(data []byte) error {
	*s = JobStatus(data)
	return nil
}
