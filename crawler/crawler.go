package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/gocolly/colly/v2"

	"github.com/leoomi/benefits-crawler/config"
)

type CrawlerInput struct {
	CPF      string
	Username string
	Password string
}

type Result struct {
	CPF      string
	Benefits []string
}

type Crawler struct {
	cfg *config.Config
}

func NewCrawler(cfg *config.Config) *Crawler {
	return &Crawler{
		cfg: cfg,
	}
}

// Using named returns just in case of panics
func (c *Crawler) GetBenefitsByCpf(in CrawlerInput) (results []Result, err error) {
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

		results = c.crawlWithRod(url, in)
	})

	coll.Visit(c.cfg.Website)

	return
}

func (c *Crawler) crawlWithRod(url string, in CrawlerInput) []Result {
	results := []Result{}
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(url).MustWindowNormal()

	ctx, cancel := context.WithCancel(context.Background())
	page = page.Context(ctx)

	go func() {
		time.Sleep(c.cfg.CrawlingTimeout)
		browser.Close()
		cancel()
	}()

	page.MustElement("#user").MustInput(in.Username)
	page.MustElement("#pass").MustInput(in.Password)
	page.MustElement("#botao").MustClick()

	page.MustWaitStable()
	page.MustElement("app-modal-fila > ion-button").MustClick()
	page.MustElement("ion-menu").MustShadowRoot().MustElement("ion-backdrop").MustClick()

	page.MustSearch("//span[text()='Encontrar Benefícios de um CPF']").MustClick()
	page.MustSearch("//ion-card-title[text()='BENEFÍCIOS DE UM CPF']")

	page.MustSearch("//ion-input").MustElement("input").MustInput(in.CPF)
	page.MustSearch("//ion-button[contains(text(), 'Procurar')]").MustClick()
	benefitEls := page.MustSearch("//ion-card-title[text()='BENEFÍCIOS ENCONTRADOS!']").
		MustParent().MustParent().
		MustElements("ion-label")

	result := Result{
		CPF:      in.CPF,
		Benefits: []string{},
	}
	for _, el := range benefitEls {
		result.Benefits = append(result.Benefits, el.MustText())
	}

	results = append(results, result)
	browser.Close()

	return results
}
