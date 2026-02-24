package main

import (
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

	log := setupLogger(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(log, ctx, cfg)
	if err != nil {
		log.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- application.Run(ctx)
	}()

	select {
	case <-ctx.Done():
	case err := <-runErrCh:
		if err != nil {
			log.Error("application run error", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("application shutdown error", "error", err)
	}
}

func setupLogger(cfg *config.Config) *slog.Logger {
	l := logger.New(logger.Config{
		Service:   cfg.LoggerConfig.Service,
		Env:       cfg.LoggerConfig.Env,
		Version:   cfg.LoggerConfig.Version,
		Level:     cfg.LoggerConfig.Level,
		AddSource: cfg.LoggerConfig.AddSource,
	})

	return l
}
