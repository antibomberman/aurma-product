package services

import (
	"aurma_product/internal/models"
	"fmt"
	"log"
	"sync"
)

var (
	productPharmacies map[string]models.ProductPharmacy
	mutex             sync.RWMutex
)

func init() {
	productPharmacies = make(map[string]models.ProductPharmacy)
}

func (s *productService) InitTotalProductPharmaciesList() error {
	limit := 1000
	offset := 0

	for {
		pharmacies, err := s.productRepository.ProductPharmacies(limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get product pharmacies: %w", err)
		}
		if len(pharmacies) == 0 {
			break
		}

		mutex.Lock()
		for _, pharmacy := range pharmacies {
			index := fmt.Sprintf("%d_%d", pharmacy.ProductId, pharmacy.PharmacyId)
			productPharmacies[index] = pharmacy
		}
		mutex.Unlock()

		offset += limit
	}

	return nil
}

func (s *productService) UpdatedProductPharmacies() ([]models.ProductElastic, error) {
	pharmacies, err := s.productRepository.ProductPharmaciesUpdated()
	if err != nil {
		return nil, fmt.Errorf("failed to get updated product pharmacies: %w", err)
	}

	var newProductPharmacies []models.ProductPharmacy
	mutex.Lock()
	for _, pharmacy := range pharmacies {
		index := fmt.Sprintf("%d_%d", pharmacy.ProductId, pharmacy.PharmacyId)
		old, exists := productPharmacies[index]
		if !exists || old.Price != pharmacy.Price || old.Count != pharmacy.Count {
			newProductPharmacies = append(newProductPharmacies, pharmacy)
		}
		productPharmacies[index] = pharmacy
	}
	mutex.Unlock()

	if len(newProductPharmacies) == 0 {
		return nil, nil
	}

	products := make([]models.ProductElastic, 0, len(newProductPharmacies))
	for _, value := range newProductPharmacies {
		product, err := s.productRepository.GetByIdSearchData(value.ProductId)
		if err != nil {
			log.Printf("Error getting product %d: %v", value.ProductId, err)
			continue
		}

		products = append(products, product.ToProductElastic(value.Price, value.Count))
	}

	return products, nil
}

func (s *productService) ProductPharmaciesList() []models.ProductElastic {
	var products []models.ProductElastic
	var wg sync.WaitGroup
	var mu sync.Mutex

	mutex.RLock()
	defer mutex.RUnlock()

	products = make([]models.ProductElastic, 0, len(productPharmacies))

	for _, value := range productPharmacies {
		wg.Add(1)
		go func(value models.ProductPharmacy) {
			defer wg.Done()

			product, err := s.productRepository.GetById(value.ProductId)
			if err != nil {
				log.Printf("Error getting product %d: %v", value.ProductId, err)
				return
			}
			mu.Lock()
			products = append(products, models.ProductElastic{
				Id:    product.Id,
				Title: product.Title,
				Price: value.Price,
				Slug:  product.Slug,
				Count: value.Count,
			})
			mu.Unlock()
		}(value)
	}

	wg.Wait()
	return products
}
