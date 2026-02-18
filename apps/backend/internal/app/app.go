package app

import (
	_ "api/docs"
	"api/internal/adapters/http"
	"api/internal/adapters/http/handlers"
	"api/internal/adapters/minio"
	"api/internal/adapters/rabbitmq"
	"api/internal/adapters/redis"
	"api/internal/services"
	"api/pkg/config"
	"context"
	"log/slog"
)

type App struct {
	s   *http.Server
	cfg *config.Config
	rdb *redis.Redis
	rmq *rabbitmq.RabbitMQ
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
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

	server, err := http.New(log, jobHandler, cfg.ServerConfig.Port)
	if err != nil {
		rmq.Close()
		return nil, err
	}

	return &App{
		s:   server,
		cfg: cfg,
		rdb: rdb,
		rmq: rmq,
	}, nil
}

func (a *App) Run() error {
	return a.s.Run()
}

func (a *App) Shutdown(ctx context.Context) {
	if a.s != nil {
		a.s.Shutdown(ctx)
	}
	if a.rmq != nil {
		a.rmq.Close()
	}
	if a.rdb != nil {
		a.rdb.Close()
	}
}
