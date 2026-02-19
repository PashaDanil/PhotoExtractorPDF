package handlers

import (
	"api/internal/domain/job"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobService interface {
	InitUpload(ctx context.Context) (*job.Job, string, error)
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
		// обработать ошибку
		return c.JSON(http.StatusInternalServerError, job.ServerErrorResponse{Error: err.Error()})
	}

	response := job.InitUploadResponse{
		JobID:     jb.JobID,
		UploadURL: uploadURL,
	}

	// обработать ошибку
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
		// обработать ошибку
		return c.JSON(
			http.StatusBadRequest, job.ServerErrorResponse{Error: "jobId is required"},
		)
	}

	err := h.jobService.CompleteUpload(c.Request().Context(), jobID)
	// обработать ошибку
	if err != nil {
		if errors.Is(err, job.ErrNotFound) {
			return c.JSON(http.StatusNotFound, job.NotFoundResponse{Error: err.Error()})
		}
		if errors.Is(err, job.ErrAlreadyCompleted) {
			return c.JSON(http.StatusConflict, job.ConflictResponse{Error: err.Error()})
		}
		if errors.Is(err, job.ErrObjectNotFound) {
			return c.JSON(http.StatusUnprocessableEntity, job.UnprocessableEntityResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, job.ServerErrorResponse{Error: err.Error()})
	}

	response := job.CompleteUploadResponse{
		JobID:  jobID,
		Status: string(job.JobStatusQueued),
	}

	// обработать ошибку
	return c.JSON(http.StatusAccepted, response)
}
