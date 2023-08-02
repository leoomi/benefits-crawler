package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/gocolly/colly/v2"

	"github.com/leoomi/benefits-crawler/config"
)

type Result struct {
	CPF      string
	Benefits []string
}

type crawler struct {
	cfg *config.Config
}

func NewCrawler(cfg *config.Config) *crawler {
	return &crawler{
		cfg: cfg,
	}
}

// Using named returns just in case of panics
func (c *crawler) GetBenefitsByCpf(cpfs []string) (results []Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			results = []Result{}
			err = fmt.Errorf("error while scraping: %v", r)
		}
	}()
	coll := colly.NewCollector()
	coll.SetRequestTimeout(c.cfg.CrawlingTimeout)

	coll.OnHTML("frame", func(e *colly.HTMLElement) {
		url := e.Attr("src")

		results = c.crawlWithRod(url, cpfs)
	})

	coll.Visit(c.cfg.Website)

	return
}

func (c *crawler) crawlWithRod(url string, cpfs []string) []Result {
	cpf := cpfs[0]
	results := []Result{}
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(url).MustWindowNormal()

	ctx, cancel := context.WithCancel(context.Background())
	page = page.Context(ctx)

	go func() {
		time.Sleep(c.cfg.CrawlingTimeout)
		cancel()
	}()

	page.MustElement("#user").MustInput("konsiteste8")
	page.MustElement("#pass").MustInput("konsiteste8")
	page.MustElement("#botao").MustClick()

	// page.MustElement("app-modal-fila > ion-button").MustClick()
	// page.MustElement("ion-menu").MustShadowRoot().MustElement("ion-backdrop").MustClick()

	// page.MustWaitStable()
	page.MustSearch("//span[text()='Encontrar Benefícios de um CPF']").MustClick()

	page.MustSearch("//ion-input").MustElement("input").MustInput(cpf)
	page.MustSearch("//ion-button[contains(text(), 'Procurar')]").MustClick()
	benefitEls := page.MustSearch("//ion-card-header").MustParent().MustElements("ion-label")

	result := Result{
		CPF:      cpf,
		Benefits: []string{},
	}
	for _, el := range benefitEls {
		result.Benefits = append(result.Benefits, el.MustText())
	}

	results = append(results, result)

	return results
}
