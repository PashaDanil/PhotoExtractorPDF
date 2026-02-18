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
	"api/internal/services"
	"api/pkg/config"
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
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	rdb, err := redis.New(cfg)
	if err != nil {
		return nil, err
	}

	mio, err := minio.New(cfg)
	if err != nil {
		_ = rdb.Close()
		return nil, err
	}

	rmq, err := rabbitmq.New(cfg)
	if err != nil {
		_ = rdb.Close()
		return nil, err
	}

	if err := rabbitmq.Setup(rmq); err != nil {
		_ = rmq.Close()
		_ = rdb.Close()
		return nil, err
	}

	jobStore := redis.NewJobStoreRepo(rdb)
	objectStorage := minio.NewObjectStorageRepo(mio)

	publisher := rabbitmq.NewPublisher(rmq.Channel())

	jobService := services.NewJobService(jobStore, objectStorage, publisher)
	jobHandler := handlers.NewJobHandler(jobService)

	server, err := rest.New(log, jobHandler, cfg.ServerConfig.Port)
	if err != nil {
		_ = rmq.Close()
		_ = rdb.Close()
		return nil, err
	}

	return &App{
		RESTserver: server,
		cfg:        cfg,
		rdb:        rdb,
		rmq:        rmq,
		pub:        publisher,
		log:        log,
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
			a.log.Error("error stopping REST server", slog.Any("err", e))
			err = errors.Join(err, e)
		}
	}

	if a.rmq != nil {
		if e := a.rmq.Close(); e != nil {
			a.log.Error("error closing rabbitmq", slog.Any("err", e))
			err = errors.Join(err, e)
		}
	}

	if a.rdb != nil {
		if e := a.rdb.Close(); e != nil {
			a.log.Error("error closing redis", slog.Any("err", e))
			err = errors.Join(err, e)
		}
	}

	return err
}
