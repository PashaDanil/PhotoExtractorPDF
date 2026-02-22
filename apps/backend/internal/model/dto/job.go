package dto

type InitResponse struct {
	JobID     string `json:"job_id"`
	UploadURL string `json:"upload_url"`
}
