package api

import (
	"log/slog"
	"microservice_t/internal/storage"

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
		storage: s,
		log:     l,
	}

	api.SetUpRoutes()

	return api
}

func (api *API) run() {

}

func (api *API) SetUpRoutes() {
	v1 := api.Router.Group("api/v1/")
	v1.GET("/url")
}
