package services

import (
	"aurma_product/internal/models"
	"aurma_product/internal/models/sadykhanModels"
	"context"
	"io"
)

// ProductService определяет интерфейс для сервиса работы с продуктами.
type ProductService interface {
	Search(query string, from, size int, sort string, minPrice, maxPrice int) ([]models.ProductDetail, int, error)
	Show(id string) (models.ProductDetail, error)
	GetImages(productID int) ([]models.ProductImage, error)
	InitTotalProductPharmaciesList() error
	UpdatedProductPharmacies() ([]models.ProductElastic, error)
	ProductPharmaciesList() []models.ProductElastic
	SetAllProductToElastic() error
}

type SadykhanService interface {
	Parse(ctx context.Context, data io.Reader, contentType string) (*sadykhanModels.Catalog, error)
	Create(ctx context.Context, catalog *sadykhanModels.Catalog) error
}
