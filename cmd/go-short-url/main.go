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
		log.Error("не удалось инициализировать хранилище", sl.Err(err))
		os.Exit(1)
	}
	id, err := storage.SaveURL("https://stepik4555.org", "stepik4555")
	if err != nil {
		log.Error("не смогли добавить ", id, sl.Err(err))
		os.Exit(1)
	}
	log.Info("без ошибок добавили")

	resStr, err := storage.GetURL("stepik")
	if err != nil {
		log.Error("не смогли дать ссылку", sl.Err(err))
		os.Exit(1)
	}
	log.Info("без ошибок получили", resStr)

	err = storage.DeleteURL("https://stepik.org")
	if err != nil {
		log.Error("не смогли удалить", sl.Err(err))
		os.Exit(1)
	}
	log.Info("без ошибок удалили")
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
