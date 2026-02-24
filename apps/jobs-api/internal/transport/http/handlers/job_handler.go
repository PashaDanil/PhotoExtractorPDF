package handlers

import (
	errs "api/internal/errors"
	"api/internal/model/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	logger     *slog.Logger
}

func NewJobHandler(jobService JobService, log *slog.Logger) *JobHandler {
	return &JobHandler{
		jobService: jobService,
		logger:     log,
	}
}

func (h *JobHandler) HandlePDFUploadRequest(c echo.Context) error {
	const op = "handlers.JobHandler.HandlePDFUploadRequest"
	ctx := c.Request().Context()

	h.logger.Debug("handling PDF upload request", "op", op)

	jb, err := h.jobService.InitUpload(ctx)
	if err != nil {
		h.logger.Error("failed to init upload", "op", op, "error", err)
		return echo.ErrInternalServerError
	}

	response := jb.ToInitResponse()

	return c.JSON(http.StatusCreated, response)
}

func (h *JobHandler) HandlePDFUploadComplete(c echo.Context) error {
	const op = "handlers.JobHandler.HandlePDFUploadComplete"
	ctx := c.Request().Context()

	jobID := c.Param("jobId")

	h.logger.Debug("handling PDF upload complete", "op", op, "job_id", jobID)

	id, err := validateJobID(jobID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = h.jobService.CompleteUpload(ctx, id); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.logger.Debug("job not found", "op", op, "job_id", id)
			return echo.NewHTTPError(http.StatusNotFound, errs.ErrNotFound.Error())
		case errors.Is(err, errs.ErrAlreadyQueued):
			h.logger.Debug("job already queued", "op", op, "job_id", id)
			return echo.NewHTTPError(http.StatusConflict, errs.ErrAlreadyQueued.Error())
		default:
			h.logger.Error("failed to complete upload", "op", op, "job_id", id, "error", err)
			return echo.ErrInternalServerError
		}
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
