package infra

import (
	"github.com/leoomi/benefits-crawler/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})

	return rdb
}
