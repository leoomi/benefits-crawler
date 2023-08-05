package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/leoomi/benefits-crawler/config"
)

const CrawlerProcessIndex = "crawler_process"
const BenefitsIndex = "benefits"

type ElsearchRes struct {
	Id string `json:"_id"`
}

type DocResponse[T any] struct {
	Id     string `json:"_id"`
	Source T      `json:"_source"`
	Found  bool   `json:"found"`
}

type UpdateReq struct {
	Doc map[string]interface{} `json:"doc"`
}

type Elsearch struct {
	client *es.Client
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

func (e *Elsearch) GetDocument(index string, id string) ([]byte, error) {
	res, err := e.client.Get(index, id)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (e *Elsearch) UpdateDocument(index string, id string, body []byte) error {
	res, err := e.client.Update(index, id, bytes.NewReader(body))
	str, _ := io.ReadAll(res.Body)
	fmt.Println(str)

	return err
}

func (e *Elsearch) SearchDocument(index string, field string, value string) ([]byte, error) {
	query := ``
	res, err := e.client.Search(
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return nil, err
	}

	str, _ := io.ReadAll(res.Body)
	return str, nil
}
