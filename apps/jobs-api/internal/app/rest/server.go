package rest

import (
	"api/internal/config"
	"api/internal/transport/http/handlers"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
}

func New(
	jobHandler *handlers.JobHandler,
	cfg *config.Config,
) *Server {
	e := echo.New()

	e.HideBanner = true
	e.Debug = false
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendConfig.URL},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	uploads := e.Group("/upload")
	uploads.POST("", jobHandler.HandlePDFUploadRequest)
	uploads.POST("/:jobId/complete", jobHandler.HandlePDFUploadComplete)

	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	const op = "rest.Server.Run"

	err := s.e.Start(":" + s.cfg.ServerConfig.Port)
	if err != nil {
		return fmt.Errorf("%s, error: %v", op, err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	const op = "rest.Server.Stop"

	err := s.e.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s, error: %v", op, err)
	}

	return nil
}
