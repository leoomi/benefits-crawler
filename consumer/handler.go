package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/leoomi/benefits-crawler/crawler"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) handleCralwerMessages(delivery amqp.Delivery) {
	fmt.Printf("Received message: %v\n", delivery)
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

		c.updateProcess(process.ID, models.Done)
		benefits := models.Benefits{
			CPF:      process.CPF,
			Benefits: result.Benefits,
		}

		err = c.createOrUpdateBenefits(benefits)
		if err != nil {
			fmt.Printf("Creating or updating ES document failed: %s\n", err)
			return
		}

		value, _ := json.Marshal(benefits)
		err = c.redis.Set(ctx, benefits.CPF, value, time.Hour).Err()
		if err != nil {
			fmt.Printf("Saving to redis failed: %s\n", err)
			return
		}
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

func (c *Consumer) createOrUpdateBenefits(benefits models.Benefits) error {
	id, err := c.elsearch.SearchSingleDocumentId(infra.BenefitsIndex, "cpf", benefits.CPF)
	if err == infra.ErrESNotFound {
		c.elsearch.CreateDocument(infra.BenefitsIndex, benefits)
		return nil
	}
	if err != nil {
		return err
	}

	c.elsearch.UpdateDocument(infra.BenefitsIndex, id, map[string]interface{}{
		"benefits": benefits.Benefits,
	})

	return nil
}
