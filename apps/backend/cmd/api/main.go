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

func main() {
	cfg := config.MustLoad()

	// log := setupLogger(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		os.Exit(1)
	}

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- application.Run(ctx)
	}()

	select {
	case <-ctx.Done():
	case err := <-runErrCh:
		_ = err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application.Shutdown(shutdownCtx)
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
