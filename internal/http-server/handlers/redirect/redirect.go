package redirect

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-short-url/internal/lib/api/response"
	"go-short-url/internal/lib/logger/sl"
	"go-short-url/internal/storage"
	"log/slog"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.redirect.New"

		log := log.With(slog.String("op", op))

		alias := c.Param("alias")
		if alias == "" {
			log.Info("empty alias")
			c.JSON(http.StatusBadRequest, response.Error("invalid alias"))
			return
		}

		url, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			c.JSON(http.StatusNotFound, response.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			c.JSON(http.StatusInternalServerError, response.Error("internal error"))
			return
		}

		log.Info("redirecting", slog.String("url", url))

		c.Redirect(http.StatusFound, url)
	}
}
