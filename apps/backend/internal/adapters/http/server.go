package http

import (
	"api/internal/adapters/http/handlers"
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	log  *slog.Logger
	e    *echo.Echo
	port string
}

func New(
	log *slog.Logger,
	jobHandler *handlers.JobHandler,
	port string,
) (*Server, error) {
	e := echo.New()

	// e.Use(middleware.RequestLogger())
	e.HideBanner = true
	e.Debug = false
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	uploads := e.Group("/uploads")
	uploads.POST("", jobHandler.HandlePDFUploadRequest)
	uploads.POST("/:jobId/complete", jobHandler.HandlePDFUploadComplete)

	return &Server{
		log:  log,
		e:    e,
		port: port,
	}, nil
}

func (s *Server) Run() error {
	const op = "http.Server.Run"

	log := s.log.With(
		slog.String("op", op),
		slog.String("port", s.port),
	)

	log.Info("starting HTTP server")

	return s.e.Start(":" + s.port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	const op = "http.Server.Shutdown"

	s.log.With(slog.String("op", op)).
		Info("shutting down HTTP server")

	return s.e.Shutdown(ctx)
}
