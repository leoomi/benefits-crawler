package crawler

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/gocolly/colly/v2"

	"github.com/leoomi/benefits-crawler/config"
)

type Result struct {
	CPF      string
	Benefits []string
}

// Using named returns just in case of panics
func GetBenefitsByCpf(config *config.Config, cpfs []string) (results []Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			results = []Result{}
			err = fmt.Errorf("error while scraping: %v", r)
		}
	}()
	c := colly.NewCollector()
	c.SetRequestTimeout(config.CrawlingTimeout)

	c.OnHTML("frame", func(e *colly.HTMLElement) {
		url := e.Attr("src")

		results = crawlWithRod(config, url, cpfs)
	})

	c.Visit(config.Website)

	return
}

func crawlWithRod(config *config.Config, url string, cpfs []string) []Result {
	cpf := cpfs[0]
	results := []Result{}
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(url).MustWindowNormal()

	page.MustElement("#user").MustInput("konsiteste8")
	page.MustElement("#pass").MustInput("konsiteste8")
	page.MustElement("#botao").MustClick()

	page.MustElement("app-modal-fila > ion-button").MustClick()
	page.MustElement("ion-menu").MustShadowRoot().MustElement("ion-backdrop").MustClick()

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
