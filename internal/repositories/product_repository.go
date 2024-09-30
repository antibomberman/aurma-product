package repositories

import (
	"aurma_product/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type productRepository struct {
	db *sqlx.DB
}

// NewProductRepository создает новый экземпляр ProductRepository.
func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

// All возвращает список всех продуктов с пагинацией.
func (r *productRepository) All(offset, limit int) ([]models.Product, error) {
	query := `SELECT id, title, slug, is_active FROM product LIMIT ? OFFSET ?`
	var products []models.Product
	err := r.db.Select(&products, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, nil
}
func (r *productRepository) AllProductSearchData(offset, limit int) ([]models.ProductSearchWithData, error) {
	query := `
			SELECT 
			product.id,
			product.title,
			product.slug,
			product.is_active,
			producers.title as company_name,
			GROUP_CONCAT(distinct barcode_products.barcode SEPARATOR ", ") AS barcode,
			GROUP_CONCAT(distinct inns.title SEPARATOR ", ") as mnn,
			GROUP_CONCAT(distinct product_form.title SEPARATOR ", ") as issue_form
			FROM product
			LEFT JOIN barcode_products ON barcode_products.product_id = product.id
			LEFT JOIN producers ON producers.id = product.producer_id
			LEFT JOIN product_inns ON product_inns.product_id = product.id
			LEFT JOIN inns ON inns.id = product_inns.inn_id
			LEFT JOIN product_product_forms  ON product_product_forms.product_id = product.id
			LEFT JOIN product_form ON product_form.id = product_product_forms.form_id
			 
			GROUP BY product.id
			
			LIMIT ? OFFSET ?
`
	var products []models.ProductSearchWithData
	err := r.db.Select(&products, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, nil
}

// GetById возвращает продукт по его ID.
func (r *productRepository) GetById(id int) (models.Product, error) {
	query := `SELECT id, title, slug FROM product WHERE id = ?`
	var product models.Product
	err := r.db.Get(&product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, fmt.Errorf("product with id %d not found", id)
		}
		return models.Product{}, fmt.Errorf("failed to fetch product: %w", err)
	}
	return product, nil
}
func (r *productRepository) GetByIdSearchData(id int) (models.ProductSearchWithData, error) {
	query := `
			SELECT 
			product.id,
			product.title,
			product.slug,
			product.is_active,
			producers.title as company_name,
			GROUP_CONCAT(distinct barcode_products.barcode SEPARATOR ", ") AS barcode,
			GROUP_CONCAT(distinct inns.title SEPARATOR ", ") as mnn,
			GROUP_CONCAT(distinct product_form.title SEPARATOR ", ") as issue_form
			FROM product
			LEFT JOIN barcode_products ON barcode_products.product_id = product.id
			LEFT JOIN producers ON producers.id = product.producer_id
			LEFT JOIN product_inns ON product_inns.product_id = product.id
			LEFT JOIN inns ON inns.id = product_inns.inn_id
			LEFT JOIN product_product_forms  ON product_product_forms.product_id = product.id
			LEFT JOIN product_form ON product_form.id = product_product_forms.form_id
			where product.id = ?
			GROUP BY product.id
`
	var product models.ProductSearchWithData
	err := r.db.Get(&product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductSearchWithData{}, fmt.Errorf("product with id %d not found", id)
		}
		return models.ProductSearchWithData{}, fmt.Errorf("failed to fetch product: %w", err)
	}
	return product, nil
}

// ProductPharmaciesUpdated возвращает список обновленных аптек с продуктами.
func (r *productRepository) ProductPharmaciesUpdated() ([]models.ProductPharmacy, error) {
	//query := `SELECT product_id, pharmacy_id, price, count FROM product_pharmacy WHERE updated_at = CURDATE()`
	query := `SELECT product_id, pharmacy_id, price, count FROM product_pharmacy`
	var productPharmacies []models.ProductPharmacy
	err := r.db.Select(&productPharmacies, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated product pharmacies: %w", err)
	}
	return productPharmacies, nil
}

// GetGalleryImages возвращает список изображений галереи для продукта.
func (r *productRepository) GetGalleryImages(productId int) ([]models.ProductGalleryImage, error) {
	query := `
		SELECT product_gallery_image.id, product_gallery_image.ownerId, product_gallery_image.name, product_gallery_image.description, product_gallery_image.rank, product_gallery_image.webp 
		FROM product_gallery_image 
		WHERE ownerId = ?
	`
	var galleryImages []models.ProductGalleryImage
	err := r.db.Select(&galleryImages, query, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gallery images for product %d: %w", productId, err)
	}
	return galleryImages, nil
}

// ProductFormName возвращает название формы продукта.
func (r *productRepository) ProductFormName(productId int) (string, error) {
	query := `
		SELECT product_form.title 
		FROM product_product_forms 
		JOIN product_form ON product_product_forms.form_id = product_form.id 
		WHERE product_product_forms.product_id = ?
	`
	var title string
	err := r.db.Get(&title, query, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no form found for product with id %d", productId)
		}
		return "", fmt.Errorf("failed to fetch product form name: %w", err)
	}
	return title, nil
}

// ProductPharmacies возвращает список аптек с продуктами с пагинацией.
func (r *productRepository) ProductPharmacies(limit, offset int) ([]models.ProductPharmacy, error) {
	query := `SELECT product_id, pharmacy_id, price, count FROM product_pharmacy LIMIT ? OFFSET ?`
	var productPharmacies []models.ProductPharmacy
	err := r.db.Select(&productPharmacies, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product pharmacies: %w", err)
	}
	return productPharmacies, nil
}

// ProductPharmacy возвращает информацию о продукте в аптеке с наименьшей ценой.
func (r *productRepository) ProductPharmacy(productId int) (models.ProductPharmacy, error) {
	query := `
		SELECT product_id, pharmacy_id, price, count 
		FROM product_pharmacy
		WHERE product_id = ?
		ORDER BY price ASC 
		LIMIT 1
	`
	var productPharmacy models.ProductPharmacy
	err := r.db.Get(&productPharmacy, query, productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductPharmacy{}, fmt.Errorf("no pharmacy found for product with id %d", productId)
		}
		return models.ProductPharmacy{}, fmt.Errorf("failed to fetch product pharmacy: %w", err)
	}
	return productPharmacy, nil
}
func (r *productRepository) ProductPharmacyBySku(sku string) (models.ProductPharmacy, error) {
	query := `
        SELECT product_id, pharmacy_id, price, count 
        FROM product_pharmacy
        WHERE barcode = ?
        ORDER BY price ASC 
        LIMIT 1
    `
	var productPharmacy models.ProductPharmacy
	err := r.db.Get(&productPharmacy, query, sku)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductPharmacy{}, fmt.Errorf("no pharmacy found for product with barcode %s", sku)
		}
		return models.ProductPharmacy{}, fmt.Errorf("failed to fetch product pharmacy by barcode: %w", err)
	}
	return productPharmacy, nil
}

func (r *productRepository) ProductUpdate(id int, updates map[string]interface{}) error {
	if _, ok := updates["updated_at"]; !ok {
		updates["updated_at"] = time.Now()
	}

	updates["id"] = id

	query := "UPDATE product SET "
	setStatements := make([]string, 0, len(updates))
	for field := range updates {
		if field != "id" {
			setStatements = append(setStatements, field+" = :"+field)
		}
	}
	query += strings.Join(setStatements, ", ")
	query += " WHERE id = :id"

	_, err := r.db.NamedExec(query, updates)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}
