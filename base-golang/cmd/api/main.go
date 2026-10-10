package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"base-golang/internal/app"
	"base-golang/internal/config"
	"base-golang/internal/platform/logger"
	"base-golang/internal/platform/postgres"
	"base-golang/migrations"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// logging
	log, err := logger.New(os.Stdout, cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		return err
	}
	slog.SetDefault(log)

	// signals
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// connect db
	startup, cancel := context.WithTimeout(ctx, 30*time.Second)
	db, err := postgres.Open(startup, cfg.DatabaseURL)
	if err != nil {
		cancel()
		return err
	}
	defer db.Close()

	// migration
	err = migrations.Apply(startup, db)
	cancel()
	if err != nil {
		return err
	}

	// run app
	return app.Run(ctx, cfg, db)
}
