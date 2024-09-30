package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"net/http"
	"time"
)

// Elastic представляет клиент Elasticsearch.
type Elastic struct {
	client *elasticsearch.Client
}

// New создает и возвращает новый экземпляр Elastic.
func New() (*Elastic, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 10,
		},
	})
	if err != nil {
		return nil, err
	}
	return &Elastic{client: es}, nil
}

// Check проверяет соединение с Elasticsearch.
func (es *Elastic) Check() error {
	ping, err := es.client.Ping()
	if err != nil {
		return err
	}
	log.Printf("Elasticsearch returned with status: %s", ping.Status())
	return nil
}
