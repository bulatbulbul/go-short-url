package main

import (
	"go-short-url/internal/config"
	"go-short-url/internal/lib/logger/sl"
	"go-short-url/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("start")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		// мб как-то полегче это можно делать
		log.Error("не удалось инициализировать хранилище", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	// TODO init router: chi

	// TODO run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log

}
