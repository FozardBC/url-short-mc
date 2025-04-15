package save

import (
	"errors"
	"log/slog"
	"microservice_t/internal/lib/api/response"
	"microservice_t/internal/lib/random"
	"microservice_t/internal/storage"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const aliasLenth = 5

type Request struct {
	URL   string `json:"url" validate:"required,url" binding:"required"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp  response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {

		log = log.With(
			slog.String("requestId:", requestid.Get(c)),
		)

		var req Request

		if err := c.BindJSON(&req); err != nil {
			log.Error("can't decode request body", "err", err.Error())

			c.JSON(http.StatusBadRequest, response.Error("failed to decode request body"))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			log.Error("invalid request", "err", err.Error())

			c.JSON(http.StatusBadRequest, response.ValidationError(validatorErr))

			return
		}

		if req.Alias == "" {
			req.Alias = random.NewRandomString(aliasLenth)
		}

		err := urlSaver.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrAliasAlreadyExists) {
				log.Debug("alias is already exists", "url", req.URL, "alias", req.Alias)

				c.JSON(http.StatusBadRequest, response.Error(err.Error()))

				return
			}

			log.Debug("can't save url", "err", err)

			c.JSON(http.StatusInternalServerError, response.Error("Internal Server Error"))
		}

		log.Debug("url saved", req.Alias, req.URL)

		//TODO: почему-то странно возвращает
		c.JSON(http.StatusOK, Response{resp: response.OK(), Alias: req.Alias})
	}

}
