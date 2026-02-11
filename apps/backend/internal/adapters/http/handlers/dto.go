package handlers

// InitUploadResponse represents the response when initializing a PDF upload
// @name InitUploadResponse
type InitUploadResponse struct {
	JobID     string `json:"job_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UploadURL string `json:"upload_url" example:"https://minio.example.com/upload?signature=..."`
}

// CompleteUploadResponse represents the response after completing upload
// @name CompleteUploadResponse
type CompleteUploadResponse struct {
	JobID  string `json:"job_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status string `json:"status" example:"queued"`
}

// ServerErrorResponse represents a server error response
// @name ServerErrorResponse
type ServerErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

// NotFoundResponse represents a not found error response
// @name NotFoundResponse
type NotFoundResponse struct {
	Error string `json:"error" example:"jobId not found"`
}

// ConflictResponse represents a conflict error response
// @name ConflictResponse
type ConflictResponse struct {
	Error string `json:"error" example:"job already completed"`
}
