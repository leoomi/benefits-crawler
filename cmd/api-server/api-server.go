package main

import (
	"flag"
	"fmt"

	"github.com/leoomi/benefits-crawler/api"
	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/infra"
	"github.com/leoomi/benefits-crawler/models"
)

func main() {
	var rootDir = flag.String("r", ".", "root directory file path to load configurations and the client files")
	flag.Parse()

	cfg, err := config.LoadConfig(*rootDir)
	if err != nil {
		panic("error reading config")
	}
	fmt.Println(*cfg)

	redis := infra.NewRedisClient(cfg)
	elsearch, err := infra.NewElsearchClient(cfg)
	if err != nil {
		fmt.Println(err)
		panic("elasticsearch connection failed")
	}

	rabbitmq, err := infra.NewAMPQClient(cfg)
	if err != nil {
		fmt.Println(err)
		panic("rabbitmq connection failed")
	}
	defer rabbitmq.Close()

	deps := models.CommondDeps{
		Cfg:      cfg,
		Redis:    redis,
		Elsearch: elsearch,
		Rabbitmq: rabbitmq,
	}

	server := api.NewServer(&deps, *rootDir)
	err = server.Start()

	if err != nil {
		fmt.Println(err)
	}
}
