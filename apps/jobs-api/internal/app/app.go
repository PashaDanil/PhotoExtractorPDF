package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"api/internal/app/rest"
	"api/internal/config"
	"api/internal/repository/cache"
	"api/internal/repository/queue"
	"api/internal/repository/storage"
	"api/internal/service"
	"api/internal/transport/http/handlers"
	"api/pkg/minio"
	"api/pkg/rabbitmq"
	"api/pkg/redis"
)

type App struct {
	logger     *slog.Logger
	RESTserver *rest.Server
	cfg        *config.Config
	rdb        *redis.Redis
	rmq        *rabbitmq.RabbitMQ
	pub        *queue.Publisher
}

func New(
	log *slog.Logger,
	ctx context.Context,
	cfg *config.Config,
) (*App, error) {
	const op = "app.New"

	rdb, err := redis.New(ctx, cfg)
	if err != nil {
		log.Error("failed to create Redis client", "op", op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("redis connected successfully")

	mio, err := minio.New(ctx, cfg)
	if err != nil {
		_ = rdb.Close()
		log.Error("failed to create MinIO client", "op", op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("minio connected successfully")

	rmq, err := rabbitmq.New(cfg)
	if err != nil {
		_ = rdb.Close()
		log.Error("failed to create RabbitMQ client", "op", op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pub := queue.NewPublisher(rmq.Channel(), log)

	jobStore := cache.NewJobStoreRepo(rdb, log)
	objectStorage := storage.NewObjectStorageRepo(mio, cfg.MinIOConfig.Bucket, log)

	jobService := service.NewJobService(jobStore, objectStorage, pub, log)
	jobHandler := handlers.NewJobHandler(jobService, log)

	RESTserver := rest.New(jobHandler, cfg)

	log.Info("app created successfully")

	return &App{
		logger:     log,
		RESTserver: RESTserver,
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
		a.logger.Info("context cancelled, shutting down")
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			a.logger.Info("REST server closed")
			return nil
		}

		a.logger.Error("REST server unexpected error", "error", err)
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
