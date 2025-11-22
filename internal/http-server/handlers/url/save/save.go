package save

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"log/slog"

	resp "go-short-url/internal/lib/api/response"
	"go-short-url/internal/lib/logger/sl"
	"go-short-url/internal/storage"
)

type Request struct {
	URL   string `json:"url" binding:"required,url"`
	Alias string `json:"alias"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(url string, alias string) (int64, error)
}

func New(log *slog.Logger, saver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.save.New"

		log := log.With(slog.String("op", op))

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to bind request", sl.Err(err))
			c.JSON(http.StatusBadRequest, resp.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = "rnd123" // временно, пока нет random
		}

		_, err := saver.SaveURL(req.URL, alias)
		if err != nil {
			if err == storage.ErrURLExists {
				log.Info("alias already exists")
				c.JSON(http.StatusConflict, resp.Error("alias already exists"))
				return
			}

			log.Error("failed to save url", sl.Err(err))
			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))
			return
		}

		c.JSON(http.StatusOK, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
