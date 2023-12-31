package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gocolly/colly/v2"

	"github.com/leoomi/benefits-crawler/config"
	"github.com/leoomi/benefits-crawler/models"
)

type CrawlerInput struct {
	CPF      string
	Username string
	Password string
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
func (c *Crawler) GetBenefitsByCpf(in CrawlerInput) (results models.Benefits, err error) {
	defer func() {
		if r := recover(); r != nil {
			results = models.Benefits{}
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

func (c *Crawler) crawlWithRod(url string, in CrawlerInput) models.Benefits {
	var page *rod.Page
	var browser *rod.Browser

	if c.cfg.RunningInContainer {
		path, _ := launcher.LookPath()
		u := launcher.New().Bin(path).MustLaunch()
		browser = rod.New().ControlURL(u)
		page = browser.MustConnect().MustPage(url)
	} else {
		browser := rod.New().MustConnect().NoDefaultDevice()
		page = browser.MustPage(url).MustWindowNormal()
	}

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

	result := models.Benefits{
		CPF:      in.CPF,
		Benefits: []string{},
	}

	if len(benefitEls) == 1 &&
		benefitEls.First().MustText() == "Matrícula não encontrada!" {
		browser.Close()
		return result
	}

	for _, el := range benefitEls {
		result.Benefits = append(result.Benefits, el.MustText())
	}

	browser.Close()

	return result
}
