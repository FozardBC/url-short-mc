package api

import (
	"errors"
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
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

func (api *API) URL() gin.HandlerFunc {
	return func(c *gin.Context) {

		api.log.With(
			"requestId", requestid.Get(c),
		)

		var req Request

		if err := c.BindJSON(&req); err != nil {
			api.log.Error("can't decode json request", "err", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			api.log.Error("invalid request", err.Error())

			c.JSON(http.StatusBadRequest, response.ValidationError(validatorErr))

			return
		}

		if req.Alias == "" {
			req.Alias = random.NewRandomString(aliasLenth)
		}

		err := api.storage.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrAliasAlreadyExists) {
				api.log.Debug("alias is already exists", "url", req.URL, "alias", req.Alias, "req_id", requestid.Get(c))
			}
		}

		c.JSON(http.StatusOK, gin.H{"msg": "url saved", "alias": req.Alias})
	}

}
