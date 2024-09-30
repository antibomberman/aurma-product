package services

import (
	"aurma_product/internal/models/sadykhanModels"
	"aurma_product/internal/repositories"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/antibomberman/dblayer"
	"io"
	"strings"
	"sync"
	"time"
)

type sadykhanService struct {
	productRepository repositories.ProductRepository
	dblayer           *dblayer.DBLayer
}

func NewSadykhanService(dblayer *dblayer.DBLayer, productRepo repositories.ProductRepository) SadykhanService {
	return &sadykhanService{productRepository: productRepo, dblayer: dblayer}
}

func (s *sadykhanService) Parse(ctx context.Context, data io.Reader, contentType string) (*sadykhanModels.Catalog, error) {
	if contentType == "" {
		return nil, fmt.Errorf("content type is empty")
	}
	content, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}

	var catalog sadykhanModels.Catalog

	contentType = strings.ToLower(contentType)
	switch strings.ToLower(contentType) {
	case "xml":
		if err := xml.Unmarshal(content, &struct {
			*sadykhanModels.Catalog
			XMLName xml.Name `xml:"rocketpharm_catalog"`
		}{Catalog: &catalog}); err != nil {
			return nil, fmt.Errorf("error parsing XML: %w", err)
		}
	case "json":
		if err := json.NewDecoder(bytes.NewReader(content)).Decode(&catalog); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return &catalog, nil
}

func (s *sadykhanService) Create(ctx context.Context, catalog *sadykhanModels.Catalog) error {

	errChan := make(chan error, len(catalog.Offers))
	wg := &sync.WaitGroup{}

	wg.Add(len(catalog.Offers))
	for _, offer := range catalog.Offers {
		go func(offer sadykhanModels.Offer) {
			defer wg.Done()
			if offer.SKU == "" {
				return
			}
			exists, err := s.dblayer.Exists(ctx, "product_pharmacy", []dblayer.Condition{{"sku", "=", offer.SKU}})

			if !exists {
				errChan <- fmt.Errorf("product with SKU %s not exists", offer.SKU)
				return
			}
			count := 0
			if offer.Availabilities.Available == "YES" {
				count = offer.Availabilities.Availability
			}
			price := 0
			if offer.CityPrices.CityPrice > 0 {
				price = offer.CityPrices.CityPrice
			}

			record, err := s.dblayer.UpdateRecord(ctx, "product_pharmacy", map[string]interface{}{
				"price":      price,
				"count":      count,
				"updated_at": time.Now(),
			}, []dblayer.Condition{{"sku", "=", offer.SKU}})
			if err != nil {
				errChan <- fmt.Errorf("error updating for SKU %s: %w", offer.SKU, err)
				return
			}
			fmt.Printf("Updated product with SKU %s: %+v\n", offer.SKU, record)

		}(offer)
	}
	wg.Wait()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during processing: %v", len(errs), errs)
	}
	return nil
}
