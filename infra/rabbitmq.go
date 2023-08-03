package infra

import (
	"context"
	"time"

	"github.com/leoomi/benefits-crawler/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

const crawlerQueueName = "crawlerProcesses"

type RMQClient struct {
	cfg  *config.Config
	conn *amqp.Connection
}

func NewAMPQClient(cfg *config.Config) (*RMQClient, error) {
	conn, err := amqp.Dial(cfg.RabbitMQAddress)
	if err != nil {
		return nil, err
	}

	declareQueue(conn)
	return &RMQClient{conn: conn}, nil
}

func declareQueue(conn *amqp.Connection) (*amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		crawlerQueueName, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	return &q, err
}

func (c *RMQClient) PublishMessage(message []byte) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"",               // exchange
		crawlerQueueName, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (c *RMQClient) GetConsumerMessages() (<-chan amqp.Delivery, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch.Consume(
		crawlerQueueName, // queue
		"",               // consumer
		true,             // auto-ack TODO change to false if time
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
}

func (c *RMQClient) Close() {
	c.conn.Close()
}
