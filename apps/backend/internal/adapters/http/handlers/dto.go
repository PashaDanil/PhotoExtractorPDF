package handlers

// InitUploadResponse represents the response when initializing a PDF upload
// @name InitUploadResponse
type InitUploadResponse struct {
	JobID     string `json:"job_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status    string `json:"status" example:"uploading"`
	UploadURL string `json:"upload_url" example:"https://minio.example.com/upload?signature=..."`
	CreatedAt int64  `json:"created_at" example:"1738800000"`
}

// CompleteUploadResponse represents the response after completing upload
// @name CompleteUploadResponse
type CompleteUploadResponse struct {
	JobID  string `json:"job_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status string `json:"status" example:"queued"`
}

// ErrorResponse represents an error response
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}
