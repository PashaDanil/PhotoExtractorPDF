package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	_ "api/docs"
	"api/internal/adapters/http/handlers"
	"api/internal/adapters/minio"
	"api/internal/adapters/rabbitmq"
	"api/internal/adapters/redis"
	"api/internal/app/rest"
	"api/internal/config"
	jobServices "api/internal/services/job"
	"log/slog"
)

type App struct {
	RESTserver *rest.Server
	cfg        *config.Config
	rdb        *redis.Redis
	rmq        *rabbitmq.RabbitMQ
	pub        *rabbitmq.Publisher
	log        *slog.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
) (*App, error) {
	rdb, err := redis.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	mio, err := minio.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	rmq, err := rabbitmq.New(cfg)
	if err != nil {
		return nil, err
	}

	pub := rabbitmq.NewPublisher(rmq.Channel())

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio, cfg.MinIOConfig.Bucket)

	jobService := jobServices.NewJobService(jobStore, objectStorage, pub)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := rest.New(jobHandler, cfg.ServerConfig.Port)
	if err != nil {
		return nil, err
	}

	return &App{
		RESTserver: server,
		cfg:        cfg,
		rdb:        rdb,
		rmq:        rmq,
		pub:        pub,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.RESTserver == nil {
		return fmt.Errorf("REST server is nil")
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- a.RESTserver.Run()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	var err error

	if a.RESTserver != nil {
		if e := a.RESTserver.Stop(ctx); e != nil {
			err = errors.Join(err, e)
		}
	}

	if a.rmq != nil {
		if e := a.rmq.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}

	if a.rdb != nil {
		if e := a.rdb.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}
