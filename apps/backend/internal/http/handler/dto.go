package handler

type InitUploadResponse struct {
	JobID        string `json:"job_id"`
	PDFKey       string `json:"pdf_key"`
	PresignedURL string `json:"presigned_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type PDFUploadSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
