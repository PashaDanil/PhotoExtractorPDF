package main

import (
	_ "api/docs"
	"api/internal/app"
	"api/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PashaDanil/logger"
)

// TODO: ошибки
// TODO: валидация http
// TODO: поднять логи
// TODO: + еще один статус в redis

// @title PDF to Images API
// @version 1.0
// @description API for converting PDF documents to images
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg)

	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Error("app init failed", slog.Any("err", err))
		os.Exit(1)
	}

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- application.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-runErrCh:
		if err != nil {
			log.Error("server stopped with error", slog.Any("err", err))
		} else {
			log.Info("server stopped")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		// обработать ошибку
	}

	log.Info("application stopped")
}

func setupLogger(cfg *config.Config) *slog.Logger {
	log := logger.New(logger.Config{
		Service:   cfg.LoggerConfig.Service,
		Env:       cfg.LoggerConfig.Env,
		Version:   cfg.LoggerConfig.Version,
		Level:     cfg.LoggerConfig.Level,
		AddSource: cfg.LoggerConfig.AddSource,
	})

	return log
}
