package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func New(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// обрабатываем запрос
		c.Next()

		latency := time.Since(start)

		log.Info("request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("latency", latency.String()),
		)
	}
}
