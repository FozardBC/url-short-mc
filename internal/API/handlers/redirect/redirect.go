package redirect

import (
	"errors"
	"log/slog"
	"microservice_t/internal/lib/api/response"
	"microservice_t/internal/storage"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.With(
			slog.String("requesId:", requestid.Get(c)),
		)

		alias := c.Param("alias")
		if alias == "" {
			log.Info("alias is empty")

			c.JSON(http.StatusBadRequest, response.Error("alias is empty"+alias))

			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrAliasNotFound) {
				log.Info("alias not found")

				c.JSON(http.StatusBadRequest, response.Error("alias not found"))

				return
			}

			log.Info("can't get url", "err", err)

			c.JSON(http.StatusInternalServerError, "Internal Server Error")

			return

		}

		c.Redirect(http.StatusPermanentRedirect, url)

		log.Info("url redirected")

		return
	}
}
