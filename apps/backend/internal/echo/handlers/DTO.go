package handlers

type ErrorResponse struct {
	Error string `json:"error"`
}

type InitUploadResponse struct {
	JobID     string `json:"jobId"`
	PDFKey    string `json:"pdfKey"`
	UploadURL string `json:"uploadUrl"`
}

type CompleteUploadResponse struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}
