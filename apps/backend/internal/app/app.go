package app

import (
	_ "api/docs"
	"api/internal/adapters/http"
	"api/internal/adapters/http/handlers"
	"api/internal/adapters/minio"
	"api/internal/adapters/redis"
	"api/internal/services"
)

type App struct {
	s *http.Server
}

func New() (*App, error) {
	a := &App{}

	rdb, err := redis.New()
	if err != nil {
		return nil, err
	}

	mio, err := minio.New()
	if err != nil {
		return nil, err
	}

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio)

	jobService := services.NewJobService(jobStore, objectStorage)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := http.New(jobHandler)
	if err != nil {
		return nil, err
	}

	a.s = server

	return a, nil
}

func (a *App) Run() {
	a.s.Run()
}
