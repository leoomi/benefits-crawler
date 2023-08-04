package consumer

import (
	"context"

	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/crawler"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Consumer struct {
	cfg      *config.Config
	rabbitmq *infra.RMQClient
	elsearch *infra.Elsearch
	redis    *redis.Client
	crawler  *crawler.Crawler
}

func NewConsumer(deps *models.CommondDeps) Consumer {
	return Consumer{
		cfg:      deps.Cfg,
		rabbitmq: deps.Rabbitmq,
		elsearch: deps.Elsearch,
		redis:    deps.Redis,
		crawler:  crawler.NewCrawler(deps.Cfg),
	}
}

func (c *Consumer) StartConsumer() error {
	msgs, err := c.rabbitmq.GetConsumerMessages()
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			c.handleCralwerMessages(d)
		}
	}()

	return nil
}
