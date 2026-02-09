package app

import (
	_ "api/docs"
	"api/internal/adapters/http"
	"api/internal/adapters/http/handlers"
	"api/internal/adapters/minio"
	"api/internal/adapters/redis"
	"api/internal/config"
	"api/internal/services"
)

type App struct {
	s   *http.Server
	cfg *config.Config
}

func New() (*App, error) {
	a := &App{}

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	a.cfg = cfg

	rdb, err := redis.New(cfg)
	if err != nil {
		return nil, err
	}

	mio, err := minio.New(cfg)
	if err != nil {
		return nil, err
	}

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio)

	jobService := services.NewJobService(jobStore, objectStorage)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := http.New(cfg, jobHandler)
	if err != nil {
		return nil, err
	}

	a.s = server

	return a, nil
}

func (a *App) Run() {
	a.s.Run()
}
