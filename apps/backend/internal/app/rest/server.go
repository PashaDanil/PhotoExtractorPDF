package rest

import (
	"api/internal/adapters/http/handlers"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	e    *echo.Echo
	port string
}

func New(
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

	uploads := e.Group("/upload")
	uploads.POST("", jobHandler.HandlePDFUploadRequest)
	uploads.POST("/:jobId/complete", jobHandler.HandlePDFUploadComplete)

	return &Server{
		e:    e,
		port: port,
	}, nil
}

func (s *Server) Run() error {
	const op = "http.Server.Run"

	err := s.e.Start(":" + s.port)
	if err != nil {
		// обработать ошибку
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	const op = "RESTserver.Stop"

	// обработать ошибку
	return s.e.Shutdown(ctx)
}
