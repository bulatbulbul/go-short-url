package main

import (
	"github.com/fatih/color"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-short-url/internal/config"
	mwLogger "go-short-url/internal/http-server/middleware/logger"
	"go-short-url/internal/lib/logger/handlers/slogpretty"
	"go-short-url/internal/lib/logger/sl"
	"go-short-url/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("start", slog.String("env", cfg.Env))
	log.Debug("debug")
	log.Error("error")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("не удалось инициализировать хранилище", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	color.NoColor = false
	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
