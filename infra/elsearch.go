package infra

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/leoomi/benefits-crawler/config"
)

const CrawlerProcessIndex = "crawler_process"
const BenefitsIndex = "benefits"

type ElsearchRes struct {
	Id string `json:"_id"`
}

type Elsearch struct {
	client *elasticsearch.Client
}

func NewElsearchClient(cfg *config.Config) (*Elsearch, error) {
	elsearchClient, err := es.NewClient(es.Config{})
	if err != nil {
		return nil, err
	}

	client := Elsearch{
		client: elsearchClient,
	}
	return &client, nil
}

func (e *Elsearch) CreateIndex(index string, data []byte) (ElsearchRes, error) {
	e.client.Indices.Create(index)
	rawElsearchRes, err := e.client.Index(index, bytes.NewReader(data))
	if err != nil {
		return ElsearchRes{}, err
	}

	bytes, err := io.ReadAll(rawElsearchRes.Body)
	if err != nil {
		return ElsearchRes{}, err
	}

	var elsearchRes ElsearchRes
	json.Unmarshal(bytes, &elsearchRes)

	return elsearchRes, nil
}
