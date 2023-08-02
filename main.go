package main

import (
	"fmt"

	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/crawler"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic("error reading config")
	}

	test, err := crawler.GetBenefitsByCpf(cfg, []string{"056.054.235-68"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", test)
}
