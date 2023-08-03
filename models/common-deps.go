package models

import (
	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/redis/go-redis/v9"
)

type CommondDeps struct {
	Cfg      *config.Config
	Redis    *redis.Client
	Elsearch *infra.Elsearch
	Rabbitmq *infra.RMQClient
}
