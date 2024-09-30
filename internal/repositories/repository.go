package repositories

import (
	"aurma_product/internal/models"
)

type ProductRepository interface {
	// GetById возвращает продукт по его ID.
	GetById(id int) (models.Product, error)

	// ProductPharmacies возвращает список аптек с продуктами с пагинацией.
	ProductPharmacies(limit, offset int) ([]models.ProductPharmacy, error)

	// ProductPharmacy возвращает информацию о продукте в аптеке с наименьшей ценой.
	ProductPharmacy(productId int) (models.ProductPharmacy, error)

	// ProductPharmaciesUpdated возвращает список обновленных аптек с продуктами.
	ProductPharmaciesUpdated() ([]models.ProductPharmacy, error)

	// GetGalleryImages возвращает список изображений галереи для продукта.
	GetGalleryImages(productId int) ([]models.ProductGalleryImage, error)

	// ProductFormName возвращает название формы продукта.
	ProductFormName(productId int) (string, error)

	// All возвращает список всех продуктов с пагинацией.
	All(offset, limit int) ([]models.Product, error)

	GetByIdSearchData(id int) (models.ProductSearchWithData, error)
	AllProductSearchData(offset, limit int) ([]models.ProductSearchWithData, error)
}
