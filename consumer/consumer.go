package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/crawler"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

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

func (c *Consumer) handleCralwerMessages(delivery amqp.Delivery) {
	go func() {
		var process models.CrawlerProcess
		json.Unmarshal(delivery.Body, &process)
		results, _ := c.crawler.GetBenefitsByCpf(crawler.CrawlerInput{
			CPF:      process.CPF,
			Username: process.Username,
			Password: process.Password,
		})

		fmt.Println(results)
	}()
}
