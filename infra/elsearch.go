package infra

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/leoomi/benefits-crawler/config"
	es "github.com/olivere/elastic/v7"
)

const CrawlerProcessIndex = "crawler_process"
const crawlerProcessMapping = `{
	"mappings": {
		"properties": {
			"cpf": { "type": "keyword" },
			"username": { "type": "text" },
			"password": { "type": "text" },
			"process_state": { "type": "text" }
		}
	}
}`
const BenefitsIndex = "benefits"
const benefitsMapping = `{
	"mappings": {
		"properties": {
			"cpf": { "type": "keyword" },
			"benefits": { "type": "text" }
		}
	}
}`

var ErrESNotFound = errors.New("document not found")

type Elsearch struct {
	client *es.Client
}

func NewElsearchClient(cfg *config.Config) (*Elsearch, error) {
	elsearchClient, err := es.NewClient()
	if err != nil {
		return nil, err
	}

	elsearch := Elsearch{
		client: elsearchClient,
	}

	err = elsearch.InitializeIndexes()
	if err != nil {
		return nil, err
	}

	return &elsearch, nil
}

func (e *Elsearch) InitializeIndexes() error {
	ctx := context.Background()
	exists, err := e.client.IndexExists(CrawlerProcessIndex).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := e.client.CreateIndex(CrawlerProcessIndex).
			BodyString(crawlerProcessMapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	exists, err = e.client.IndexExists(BenefitsIndex).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		_, err := e.client.CreateIndex(BenefitsIndex).
			BodyString(benefitsMapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Elsearch) GetDocument(index string, id string, doc interface{}) error {
	get, err := e.client.
		Get().
		Index(index).
		Id(id).
		Do(context.Background())

	if err != nil {
		return err
	}

	if !get.Found {
		doc = nil
		return ErrESNotFound
	}

	fields, _ := json.Marshal(get.Source)
	json.Unmarshal(fields, doc)

	return nil
}

func (e *Elsearch) CreateDocument(index string, doc interface{}) (string, error) {
	put1, err := e.client.Index().
		Index(index).
		BodyJson(doc).
		Do(context.Background())

	if err != nil {
		return "", err
	}

	return put1.Id, nil
}

func (e *Elsearch) UpdateDocument(index string, id string, doc map[string]interface{}) error {
	_, err := e.client.
		Update().
		Index(index).
		Id(id).
		Doc(doc).
		Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (e *Elsearch) SearchSingleDocument(index string, field string, value string, doc interface{}) error {
	query := es.NewTermQuery(field, value)
	res, err := e.client.
		Search().
		Index(index).
		Query(query).
		Do(context.Background())

	if err != nil {
		return err
	}

	if res.Hits.TotalHits.Value == 0 {
		return ErrESNotFound
	}

	for _, h := range res.Hits.Hits {
		source, _ := h.Source.MarshalJSON()
		json.Unmarshal(source, doc)
	}

	return nil
}

func (e *Elsearch) SearchSingleDocumentId(index string, field string, value string) (string, error) {
	query := es.NewTermQuery(field, value)
	res, err := e.client.
		Search().
		Index(index).
		Query(query).
		Do(context.Background())

	if err != nil {
		return "", err
	}

	if res.Hits.TotalHits.Value == 0 {
		return "", ErrESNotFound
	}

	for _, h := range res.Hits.Hits {
		return h.Id, nil
	}

	return "", errors.New("something went wrong")
}
