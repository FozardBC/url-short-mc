package delete

import (
	"errors"
	"log/slog"
	"microservice_t/internal/lib/api/response"
	"microservice_t/internal/storage"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {

		log = log.With(
			slog.String("requestId:", requestid.Get(c)),
		)

		alias := c.Param("alias")

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrAliasNotFound) {
				log.Error("alias not found", "alias", alias)

				c.JSON(http.StatusBadRequest, response.Error("Alias not exists"))

				return
			}
			log.Error("can't delete url", "err", err.Error())

			c.JSON(http.StatusInternalServerError, response.Error("Internal server error"))

			return
		}

		log.Debug("url deleted", "alias", alias)

		c.JSON(http.StatusOK, "alias deleted")
	}
}
