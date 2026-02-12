package app

import (
	_ "api/docs"
	"api/internal/adapters/http"
	"api/internal/adapters/http/handlers"
	"api/internal/adapters/minio"
	"api/internal/adapters/rabbitmq"
	"api/internal/adapters/redis"
	"api/internal/config"
	"api/internal/services"
)

type App struct {
	s   *http.Server
	cfg *config.Config
	rdb *redis.Redis
	rmq *rabbitmq.RabbitMQ
}

func New() (*App, error) {
	a := &App{}

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

	rmq, err := rabbitmq.New(cfg)
	if err != nil {
		return nil, err
	}

	if err := rabbitmq.Setup(rmq); err != nil {
		rmq.Close()
		return nil, err
	}

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio)
	publisher := rabbitmq.NewPublisher(rmq.Channel())

	jobService := services.NewJobService(jobStore, objectStorage, publisher)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := http.New(cfg, jobHandler)
	if err != nil {
		rmq.Close()
		return nil, err
	}

	a.s = server
	a.rdb = rdb
	a.rmq = rmq

	return a, nil
}

func (a *App) Run() {
	a.s.Run()
}

func (a *App) Shutdown() {
	if a.s != nil {
		a.s.Shutdown()
	}
	if a.rmq != nil {
		a.rmq.Close()
	}
	if a.rdb != nil {
		a.rdb.Close()
	}
}
