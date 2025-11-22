package url

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func NewDelete(log *slog.Logger, storage URLDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		alias := c.Param("alias")
		if alias == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "alias is required",
			})
			return
		}

		err := storage.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", "alias", alias, "error", err.Error())
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  "alias not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"alias":  alias,
		})
	}
}
