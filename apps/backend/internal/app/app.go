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
	"api/internal/services"
)

type App struct {
	RESTserver *rest.Server
	cfg        *config.Config
	rdb        *redis.Redis
	rmq        *rabbitmq.RabbitMQ
	pub        *rabbitmq.Publisher
}

func New(
	ctx context.Context,
	cfg *config.Config,
) (*App, error) {
	rdb, err := redis.New(ctx, cfg)
	if err != nil {
		// обработать ошибку
		return nil, err
	}

	mio, err := minio.New(ctx, cfg)
	if err != nil {
		// обработать ошибку
		_ = rdb.Close()

		return nil, err
	}

	rmq, err := rabbitmq.New(cfg)
	if err != nil {
		// обработать ошибку
		return nil, err
	}

	if err := rabbitmq.Setup(rmq); err != nil {
		_ = rmq.Close()
		// обработать ошибку
		return nil, err
	}

	pub := rabbitmq.NewPublisher(rmq.Channel())

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio, cfg.MinIOConfig.Bucket)

	jobService := services.NewJobService(jobStore, objectStorage, pub)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := rest.New(jobHandler, cfg.ServerConfig.Port)
	if err != nil {
		_ = rmq.Close()
		_ = rdb.Close()
		// обработать ошибку
		return nil, err
	}

	return &App{
		RESTserver: server,
		cfg:        cfg,
		rdb:        rdb,
		rmq:        rmq,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.RESTserver == nil {
		// обработать ошибку
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
		// обработать ошибку
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		// обработать ошибку
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	var err error
	// обработать ошибку

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
