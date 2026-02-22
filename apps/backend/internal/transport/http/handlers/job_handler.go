package handlers

import (
	"api/internal/domain"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JobService interface {
	InitUpload(ctx context.Context) (*domain.Job, error)
	CompleteUpload(ctx context.Context, jobID uuid.UUID) error
}

type JobHandler struct {
	jobService JobService
}

func NewJobHandler(jobService JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	ctx := c.Request().Context()

	jb, err := h.jobService.InitUpload(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}

	response := jb.ToInitResponse()

	return c.JSON(http.StatusCreated, response)
}

func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	ctx := c.Request().Context()

	jobID := c.Param("jobId")

	id, err := validateJobID(jobID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.jobService.CompleteUpload(ctx, id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"job_id": id.String(),
	})
}

func validateJobID(jobID string) (uuid.UUID, error) {
	if jobID == "" {
		return uuid.Nil, fmt.Errorf("jobId is required")
	}

	id, err := uuid.Parse(jobID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid jobID: %w", err)
	}

	return id, nil
}
