package app

import (
	"context"
	"go-api/internal/echo"
	"go-api/internal/echo/handlers"
	"go-api/internal/redis"
	"go-api/internal/service"
	"go-api/internal/storage"
	"log"

	"github.com/joho/godotenv"
)

type App struct {
	s *echo.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	_ = godotenv.Load(".env")

	rdb, err := redis.New(ctx)
	if err != nil {
		return nil, err
	}

	minioClient, err := storage.New(ctx)
	if err != nil {
		return nil, err
	}

	jobRedis := redis.NewJobRedis(rdb)
	minioStorage := storage.NewStorage(minioClient)

	jobService := service.NewJobService(jobRedis, minioStorage)
	jobHandler := handlers.NewJobHandler(jobService)

	s, err := echo.New(
		jobHandler,
	)
	if err != nil {
		return nil, err
	}

	a.s = s

	log.Println("app initialized")

	return a, nil
}

func (a *App) Run() error {
	a.s.Run()
	log.Println("server started on :8080")

	return nil
}
