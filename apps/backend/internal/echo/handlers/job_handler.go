package handlers

import (
	"go-api/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
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
// @Summary Initialize PDF upload
// @Description Create a new job and get a presigned URL for uploading a PDF file
// @Tags Jobs
// @Accept json
// @Produce json
// @Success 200 {object} InitUploadResponse "Upload initialized successfully"
// @Failure 500 {object} ErrorResponse "Failed to initialize upload"
// @Router /upload [post]
func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	job, err := h.jobService.InitUpload(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, InitUploadResponse{
		JobID:     job.JobID,
		PDFKey:    job.PDFKey,
		UploadURL: job.UploadURL,
	})
}

// HandlePDFUploadComplete godoc
// @Summary Complete PDF upload
// @Description Mark the PDF upload as complete and start processing
// @Tags Jobs
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 200 {object} CompleteUploadResponse "Upload completed successfully"
// @Failure 500 {object} ErrorResponse "Failed to complete upload"
// @Router /upload/{jobId}/complete [post]
func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	jobID := c.Param("jobId")

	err := h.jobService.CompleteUpload(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, CompleteUploadResponse{
		JobID:  jobID,
		Status: "queued",
	})
}
