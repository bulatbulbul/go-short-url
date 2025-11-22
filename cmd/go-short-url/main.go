package main

import (
	"github.com/gin-gonic/gin"
	url "go-short-url/internal/http-server/handlers/url/delete"
	"log/slog"
	"net/http"
	"os"

	"go-short-url/internal/config"
	"go-short-url/internal/http-server/handlers/redirect"
	"go-short-url/internal/http-server/handlers/url/save"
	mwLogger "go-short-url/internal/http-server/middleware/logger"
	"go-short-url/internal/lib/logger/handlers/slogpretty"
	"go-short-url/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting", slog.String("env", cfg.Env))

	// Storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		os.Exit(1)
	}

	// Gin
	router := gin.New()
	router.Use(mwLogger.New(log))
	router.Use(gin.Recovery())

	// POST /add
	router.POST("/add", save.New(log, storage))

	// GET /:alias
	router.GET("/:alias", redirect.New(log, storage))

	deleteHandler := url.NewDelete(log, storage)
	router.DELETE("/delete/:alias", deleteHandler)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server running", slog.String("address", cfg.Address))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error", slog.Any("err", err))
	}
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case "local":
		handler := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
		}.NewPrettyHandler(os.Stdout)
		return slog.New(handler)

	case "dev":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return slog.Default()
}
