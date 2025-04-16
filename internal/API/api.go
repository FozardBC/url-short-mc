package api

import (
	"log/slog"
	"microservice_t/internal/API/handlers/redirect"
	deleteHandler "microservice_t/internal/API/handlers/url/delete"
	"microservice_t/internal/API/handlers/url/save"
	"microservice_t/internal/storage"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type API struct {
	Router  *gin.Engine
	storage storage.Storage
	log     *slog.Logger
}

func NewAPI(l *slog.Logger, s storage.Storage) *API {
	api := &API{
		Router:  gin.New(),
		Router:  gin.New(),
		storage: s,
		log:     l,
	}

	api.SetUpRoutes()

	return api
}

func (api *API) SetUpRoutes() {
	v1 := api.Router.Group("api/v1/")
	v1.Use(requestid.New())
	v1.Use(gin.Logger())

	v1.POST("/url", save.New(api.log, api.storage))
	v1.DELETE("/url/:alias", deleteHandler.New(api.log, api.storage))
	v1.Use(gin.Logger())

	v1.POST("/url", save.New(api.log, api.storage))
	v1.DELETE("/url/:alias", deleteHandler.New(api.log, api.storage))

	v1.GET("/:alias", redirect.New(api.log, api.storage))
	v1.GET("/:alias", redirect.New(api.log, api.storage))

}
