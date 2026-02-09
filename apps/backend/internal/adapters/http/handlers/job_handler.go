package handlers

import (
	"api/internal/domain/job"
	"api/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

// HandlePDFUploadRequest godoc
// @Summary Initialize PDF upload
// @Description Initialize a new PDF upload job and get presigned upload URL
// @Tags upload
// @Accept json
// @Produce json
// @Success 200 {object} InitUploadResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload [post]
func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	job, err := h.jobService.InitUpload(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	response := InitUploadResponse{
		JobID:     job.JobID,
		Status:    string(job.Status),
		UploadURL: job.UploadURL,
		CreatedAt: job.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// HandlePDFUploadComplete godoc
// @Summary Complete PDF upload
// @Description Mark the PDF upload as complete and queue for processing
// @Tags upload
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 200 {object} CompleteUploadResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload/{jobId}/complete [post]
func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	jobID := c.Param("jobId")

	err := h.jobService.CompleteUpload(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	response := CompleteUploadResponse{
		JobID:  jobID,
		Status: string(job.JobStatusQueued),
	}

	return c.JSON(http.StatusOK, response)
}
