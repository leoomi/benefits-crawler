package api

import (
	"github.com/gin-gonic/gin"
	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
	"github.com/redis/go-redis/v9"
)

type server struct {
	cfg      *config.Config
	router   *gin.Engine
	redis    *redis.Client
	elsearch *infra.Elsearch
	rabbitmq *infra.RMQClient
}

func NewServer(deps *models.CommondDeps) *server {
	api := &server{
		cfg:      deps.Cfg,
		redis:    deps.Redis,
		elsearch: deps.Elsearch,
		rabbitmq: deps.Rabbitmq,
	}

	api.setupRoutes()

	return api
}

func (a *server) Start() error {
	return a.router.Run(a.cfg.ServerAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
