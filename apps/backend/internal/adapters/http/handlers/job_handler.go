package handlers

import (
	"api/internal/domain/job"
	"api/internal/services"
	"api/pkg/errorx"
	"errors"
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
// @Success 201 {object} InitUploadResponse
// @Failure 500 {object} ServerErrorResponse
// @Router /upload [post]
func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	jb, err := h.jobService.InitUpload(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ServerErrorResponse{Error: err.Error()})
	}

	response := InitUploadResponse{
		JobID:     jb.JobID,
		UploadURL: jb.UploadURL,
	}

	return c.JSON(http.StatusCreated, response)
}

// HandlePDFUploadComplete godoc
// @Summary Complete PDF upload
// @Description Mark the PDF upload as complete and queue for processing
// @Tags upload
// @Accept json
// @Produce json
// @Param jobId path string true "Job ID"
// @Success 202 {object} CompleteUploadResponse
// @Failure 404 {object} NotFoundResponse
// @Failure 409 {object} ConflictResponse
// @Failure 500 {object} ServerErrorResponse
// @Router /upload/{jobId}/complete [post]
func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	jobID := c.Param("jobId")

	err := h.jobService.CompleteUpload(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return c.JSON(http.StatusNotFound, NotFoundResponse{Error: err.Error()})
		}
		if errors.Is(err, errorx.ErrAlreadyCompleted) {
			return c.JSON(http.StatusConflict, ConflictResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, ServerErrorResponse{Error: err.Error()})
	}

	response := CompleteUploadResponse{
		JobID:  jobID,
		Status: string(job.JobStatusQueued),
	}

	return c.JSON(http.StatusAccepted, response)
}
