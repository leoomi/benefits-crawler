package main

import (
	"flag"
	"fmt"

	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/consumer"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
)

func main() {
	var configFlag = flag.String("c", ".", "config file path")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFlag)
	if err != nil {
		panic("error reading config")
	}

	redis := infra.NewRedisClient(cfg)
	elsearch, err := infra.NewElsearchClient(cfg)
	if err != nil {
		panic("elasticsearch connection failed")
	}

	rabbitmq, err := infra.NewAMPQClient(cfg)
	if err != nil {
		panic("rabbitmq connection failed")
	}
	defer rabbitmq.Close()

	deps := models.CommondDeps{
		Cfg:      cfg,
		Redis:    redis,
		Elsearch: elsearch,
		Rabbitmq: rabbitmq,
	}

	var forever chan struct{}
	cons := consumer.NewConsumer(&deps)
	err = cons.StartConsumer()
	if err != nil {
		fmt.Println("error starting the consumer: ", err)
	}

	fmt.Println("Crawler consumer is listening to messages")
	<-forever
}
