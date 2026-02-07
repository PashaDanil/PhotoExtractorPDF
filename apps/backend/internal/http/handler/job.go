package handler

import (
	"encoding/json"
	"go-api/internal/service"
	"log"
	"net/http"
	"strings"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

// HandlePDFUploadRequest godoc
// @ID initPDFUpload
// @Summary Initialize PDF upload
// @Description Create a new job and get a presigned URL for uploading a PDF file
// @Tags Jobs
// @Accept json
// @Produce json
// @Success 201 {object} InitUploadResponse "Upload initialized successfully"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Failure 500 {object} ErrorResponse "Failed to initialize upload"
// @Router /jobs [post]
func (h *JobHandler) HandlePDFUploadRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	job, err := h.jobService.InitUpload(r.Context())
	if err != nil {
		log.Printf("init upload failed: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to initialize upload")
		return
	}

	writeJSON(w, http.StatusCreated, InitUploadResponse{
		JobID:        job.JobID,
		PDFKey:       job.PDFKey,
		PresignedURL: job.UploadURL,
	})
}

// HandlePDFUploadComplete godoc
// @ID completePDFUpload
// @Summary Complete PDF upload
// @Description Mark the PDF upload as complete and start processing
// @Tags Jobs
// @Accept json
// @Produce json
// @Param job_id path string true "Job ID"
// @Success 200 {object} SuccessResponse "Upload completed successfully"
// @Failure 400 {object} ErrorResponse "Invalid path or job id"
// @Failure 405 {object} ErrorResponse "Method not allowed"
// @Failure 500 {object} ErrorResponse "Failed to complete upload"
// @Router /jobs/{job_id}/complete-upload [post]
func (h *JobHandler) HandlePDFUploadComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	jobID, ok := parseCompleteUploadPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid path or job id")
		return
	}

	if err := h.jobService.CompleteUpload(r.Context(), jobID); err != nil {
		log.Printf("complete upload failed job_id=%s err=%v", jobID, err)
		writeError(w, http.StatusInternalServerError, "Failed to complete upload")
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Upload completed successfully",
	})
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func parseCompleteUploadPath(p string) (string, bool) {
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) != 3 || parts[0] != "jobs" || parts[2] != "complete-upload" {
		return "", false
	}
	jobID := parts[1]
	if jobID == "" || strings.Contains(jobID, "/") {
		return "", false
	}
	return jobID, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
