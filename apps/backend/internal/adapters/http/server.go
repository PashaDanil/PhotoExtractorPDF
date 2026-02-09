package http

import (
	"api/internal/adapters/http/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	e *echo.Echo
}

func New(
	jobHandler *handlers.JobHandler,
) (*Server, error) {
	s := &Server{}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/upload", jobHandler.HandlePDFUploadRequest)
	e.POST("/upload/:jobId/complete", jobHandler.HandlePDFUploadComplete)

	s.e = e

	return s, nil
}

func (s *Server) Run() {
	s.e.Start(":8080")
}
