package handlers

import (
	"api/internal/domain"
	"api/internal/errorx"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobService interface {
	InitUpload(ctx context.Context) (*domain.Job, string, error)
	CompleteUpload(ctx context.Context, jobID string) error
}

type JobHandler struct {
	jobService JobService
}

func NewJobHandler(jobService JobService) *JobHandler {
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
// @Success 201 {object} job.InitUploadResponse
// @Failure 500 {object} job.ServerErrorResponse
// @Router /upload [post]
func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	jb, uploadURL, err := h.jobService.InitUpload(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ServerErrorResponse{Error: err.Error()})
	}

	response := domain.InitUploadResponse{
		JobID:     jb.JobID,
		UploadURL: uploadURL,
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
// @Success 202 {object} job.CompleteUploadResponse
// @Failure 404 {object} job.NotFoundResponse
// @Failure 409 {object} job.ConflictResponse
// @Failure 422 {object} job.UnprocessableEntityResponse
// @Failure 500 {object} job.ServerErrorResponse
// @Router /upload/{jobId}/complete [post]
func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	jobID := c.Param("jobId")

	if jobID == "" {
		return c.JSON(
			http.StatusBadRequest, domain.ServerErrorResponse{Error: "jobId is required"},
		)
	}

	err := h.jobService.CompleteUpload(c.Request().Context(), jobID)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return c.JSON(http.StatusNotFound, domain.NotFoundResponse{Error: err.Error()})
		}
		if errors.Is(err, errorx.ErrAlreadyCompleted) {
			return c.JSON(http.StatusConflict, domain.ConflictResponse{Error: err.Error()})
		}
		if errors.Is(err, errorx.ErrObjectNotFound) {
			return c.JSON(http.StatusUnprocessableEntity, domain.UnprocessableEntityResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, domain.ServerErrorResponse{Error: err.Error()})
	}

	response := domain.CompleteUploadResponse{
		JobID:  jobID,
		Status: string(domain.JobStatusQueued),
	}

	return c.JSON(http.StatusAccepted, response)
}
