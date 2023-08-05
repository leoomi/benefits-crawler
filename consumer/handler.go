package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/leoomi/benefits-crawler/crawler"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) handleCralwerMessages(delivery amqp.Delivery) {
	go func() {
		var process models.CrawlerProcessWithId
		json.Unmarshal(delivery.Body, &process)

		if c.isCPFInCache(process.CPF) {
			c.updateProcess(process.ID, models.Canceled)
			return
		}

		c.updateProcess(process.ID, models.Running)
		result, err := c.crawler.GetBenefitsByCpf(crawler.CrawlerInput{
			CPF:      process.CPF,
			Username: process.Username,
			Password: process.Password,
		})

		if err != nil {
			fmt.Printf("Crawler failed: %s\n", err)
			c.updateProcess(process.ID, models.Failed)
			return
		}

		value, _ := json.Marshal(result)
		err = c.redis.Set(ctx, result.CPF, value, -1).Err()
		if err != nil {
			fmt.Printf("Saving to redis failed: %s\n", err)
		}

		c.updateProcess(process.ID, models.Done)

		benefits := models.Benefits{
			CPF:      process.CPF,
			Benefits: result.Benefits,
		}

		c.redis.Set(ctx, benefits.CPF, benefits, 0)
		c.elsearch.CreateDocument(infra.BenefitsIndex, benefits)
	}()
}

func (c *Consumer) isCPFInCache(cpf string) bool {
	_, err := c.redis.Get(ctx, cpf).Result()

	return err == nil
}

func (c *Consumer) updateProcess(id string, state models.ProcessState) error {
	doc := map[string]interface{}{
		"process_state": state,
	}
	err := c.elsearch.UpdateDocument(infra.CrawlerProcessIndex, id, doc)
	return err
}
