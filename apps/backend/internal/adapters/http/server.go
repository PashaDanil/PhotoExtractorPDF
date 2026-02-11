package http

import (
	"api/internal/adapters/http/handlers"
	"api/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
}

func New(
	cfg *config.Config,
	jobHandler *handlers.JobHandler,
) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}

	e := echo.New()

	// e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/upload", jobHandler.HandlePDFUploadRequest)
	e.POST("/upload/:jobId/complete", jobHandler.HandlePDFUploadComplete)

	s.e = e

	return s, nil
}

func (s *Server) Run() {
	s.e.Start(":" + s.cfg.Server.Port)
}
