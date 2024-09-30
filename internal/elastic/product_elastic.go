package elastic

import (
	"aurma_product/internal/models/elasticModels"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strings"
)

const IndexName = "product_list"

// ProductCreateIndex создает индекс для продуктов в Elasticsearch.
func (es *Elastic) ProductCreateIndex() error {
	settings := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"my_analyzer": map[string]interface{}{
						"tokenizer":   "whitespace",
						"char_filter": []string{"e_mapping", "rus_en_key"},
						"filter":      []string{"lowercase", "complex_word_decompound"},
					},
					"ngram_analyzer": map[string]interface{}{
						"tokenizer": "ngram",
						"filter":    []string{"lowercase"},
					},
				},
				"tokenizer": map[string]interface{}{
					"ngram": map[string]interface{}{
						"type":     "edge_ngram",
						"min_gram": 3,
						"max_gram": 10,
					},
				},
				"char_filter": map[string]interface{}{
					"e_mapping": map[string]interface{}{
						"type":     "mapping",
						"mappings": []string{"Ё=>Е", "ё=>е"},
					},
					"rus_en_key": map[string]interface{}{
						"type": "mapping",
						"mappings": []string{
							"a => ф", "b => и", "c => с", "d => в", "e => у", "f => а", "g => п", "h => р", "i => ш",
							"j => о", "k => л", "l => д", "m => ь", "n => т", "o => щ", "p => з", "r => к", "s => ы",
							"t => е", "u => г", "v => м", "w => ц", "x => ч", "y => н", "z => я", "A => Ф", "B => И",
							"C => С", "D => В", "E => У", "F => А", "G => П", "H => Р", "I => Ш", "J => О", "K => Л",
							"L => Д", "M => Ь", "N => Т", "O => Щ", "P => З", "R => К", "S => Ы", "T => Е", "U => Г",
							"V => М", "W => Ц", "X => Ч", "Y => Н", "Z => Я", "[ => х", "] => ъ", "; => ж", "< => б",
							"> => ю", ", => б", ". => ю",
						},
					},
				},
				"filter": map[string]interface{}{
					"complex_word_decompound": map[string]interface{}{
						"type":             "dictionary_decompounder",
						"word_list":        []string{"аква", "марис"},
						"max_subword_size": 22,
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":     "text",
					"analyzer": "my_analyzer",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{
							"type": "keyword",
						},
						"ngram": map[string]interface{}{
							"type":     "text",
							"analyzer": "ngram_analyzer",
						},
					},
				},
				"company_name": map[string]interface{}{
					"type":     "text",
					"analyzer": "my_analyzer",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{
							"type": "keyword",
						},
						"ngram": map[string]interface{}{
							"type":     "text",
							"analyzer": "ngram_analyzer",
						},
					},
				},
				"barcode": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"mnn": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"issue_form": map[string]interface{}{
					"type": "text",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(settings); err != nil {
		return fmt.Errorf("error encoding index settings: %w", err)
	}

	res, err := es.client.Indices.Create(
		IndexName,
		es.client.Indices.Create.WithBody(&buf),
		es.client.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		}
		return fmt.Errorf("error creating index: %v", e)
	}

	return nil
}

// ProductAddDocument добавляет продукты в Elasticsearch.
func (es *Elastic) ProductAddDocument(products []elasticModels.Product) error {
	var buf bytes.Buffer
	for _, product := range products {
		meta := []byte(fmt.Sprintf(`{"index":{"_index":"%s","_id":"%d"}}%s`, IndexName, product.Id, "\n"))
		data, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product: %w", err)
		}
		buf.Grow(len(meta) + len(data) + 1)
		buf.Write(meta)
		buf.Write(data)
		buf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Index:   IndexName,
		Body:    &buf,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("failed to perform bulk request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("failed to parse error response: %w", err)
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

// ProductSearch выполняет поиск продуктов по тексту.
func (es *Elastic) ProductSearch(ctx context.Context, text string, from, size int, sort string, minPrice, maxPrice int) ([]elasticModels.Product, int, error) {
	boolQuery := map[string]interface{}{
		"must": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":                text,
				"type":                 "best_fields",
				"fields":               []string{"name", "company_name", "barcode", "mnn", "issue_form"},
				"minimum_should_match": "90%",
				"fuzziness":            "AUTO:0,2",
			},
		},
	}

	if minPrice >= 0 || maxPrice > 0 {
		fmt.Println("minPrice >= 0 || maxPrice > 0")
		priceRange := map[string]interface{}{}
		if minPrice > 0 {
			fmt.Println("minPrice > 0")
			priceRange["gte"] = minPrice
		}
		if maxPrice > 0 {
			fmt.Println("maxPrice > 0")
			priceRange["lte"] = maxPrice
		}
		boolQuery["filter"] = []map[string]interface{}{
			{
				"range": map[string]interface{}{
					"price": priceRange,
				},
			},
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"_source": true,
		"from":    from,
		"size":    size,
	}

	// Добавление сортировки
	switch strings.ToUpper(sort) {
	case "PRICE_DESC":
		query["sort"] = []map[string]interface{}{
			{"price": map[string]interface{}{"order": "desc"}},
		}
	case "PRICE_ASC":
		query["sort"] = []map[string]interface{}{
			{"price": map[string]interface{}{"order": "asc"}},
		}
	case "COUNT_DESC":
		query["sort"] = []map[string]interface{}{
			{"count": map[string]interface{}{"order": "desc"}},
		}
	case "COUNT_ASC":
		query["sort"] = []map[string]interface{}{
			{"count": map[string]interface{}{"order": "asc"}},
		}
	case "DEFAULT", "":
	default:
		return nil, 0, fmt.Errorf("unknown sort option: %s", sort)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(IndexName),
		es.client.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to perform search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source elasticModels.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}
	total := result.Hits.Total.Value
	products := make([]elasticModels.Product, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		products[i] = hit.Source
	}

	return products, total, nil
}

// ProductSearchIds выполняет поиск ID продуктов по тексту.
func (es *Elastic) ProductSearchIds(text string, from, size int) ([]int, error) {
	//query := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"multi_match": map[string]interface{}{
	//			"query":     text,
	//			"fields":    []string{"name", "company_name", "barcode", "mnn", "issue_form"},
	//			"type":      "best_fields",
	//			"fuzziness": "AUTO:0,2",
	//		},
	//	},
	//	"_source": false,
	//	"from":    from * size,
	//	"size":    size,
	//}
	query := `{
	  "bool": {
		"should": [
		  {
			"multi_match": {
			  "query": "` + text + `",
			  "type": "best_fields",
			  "fields": ["name", "company_name", "barcode", "mnn", "issue_form"],
			  "minimum_should_match": "90%",
			  "fuzziness": "AUTO:0,2"
			}
		  }
		]
	  }
	}`

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(IndexName),
		es.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				ID int `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	ids := make([]int, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		ids[i] = hit.ID
	}

	return ids, nil
}
