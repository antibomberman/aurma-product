package services

import (
	"aurma_product/internal/models"
	"aurma_product/internal/models/elasticModels"
	"aurma_product/internal/models/sadykhanModels"
	"context"
	"io"
)

// ProductService определяет интерфейс для сервиса работы с продуктами.
type ProductService interface {
	Search(ctx context.Context, query string, from, size int, sort string, minPrice, maxPrice int) ([]models.ProductDetail, int, error)
	Show(ctx context.Context, id string) (models.ProductDetail, error)
	GetImages(productID int) ([]models.ProductImage, error)
	InitTotalProductPharmaciesList() error
	UpdatedProductPharmacies() ([]elasticModels.Product, error)
	ProductPharmaciesList() []elasticModels.Product
	SetAllProductToElastic() error
}

type SadykhanService interface {
	Parse(ctx context.Context, data io.Reader, contentType string) (*sadykhanModels.Catalog, error)
	Create(ctx context.Context, catalog *sadykhanModels.Catalog) error
}
