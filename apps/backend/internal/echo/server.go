package echo

import (
	"go-api/internal/echo/handlers"

	"github.com/labstack/echo/v4"
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

	e.POST("/upload", jobHandler.HandlePDFUploadRequest)
	e.POST("/upload/:jobId/complete", jobHandler.HandlePDFUploadComplete)
	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	s.e = e

	return s, nil
}

func (s *Server) Run() {
	s.e.Start(":8080")
}
