package services

import (
	"aurma_product/internal/elastic"
	"aurma_product/internal/models"
	"aurma_product/internal/repositories"
	"fmt"
	"github.com/antibomberman/dblayer"
	"log"
)

type productService struct {
	productRepository repositories.ProductRepository
	dblayer           *dblayer.DBLayer
	elastic           *elastic.Elastic
}

func NewProductService(dblayer *dblayer.DBLayer, productRepo repositories.ProductRepository, es *elastic.Elastic) ProductService {
	return &productService{productRepository: productRepo, dblayer: dblayer, elastic: es}
}

func (s *productService) Search(text string, from, size int, sort string, minPrice, maxPrice int) ([]models.ProductDetail, int, error) {
	elasticProducts, total, err := s.elastic.ProductSearch(text, from, size, sort, minPrice, maxPrice)
	if err != nil {
		return nil, total, fmt.Errorf("failed to search product IDs: %w", err)
	}
	productDetails := make([]models.ProductDetail, len(elasticProducts))

	for i, elasticProduct := range elasticProducts {
		images, _ := s.GetImages(elasticProduct.Id)
		productDetails[i] = elasticProduct.ToProductDetail(images)
	}

	return productDetails, total, nil

}

func (s *productService) Show(id string) (models.ProductDetail, error) {
	// Реализация метода Show
	return models.ProductDetail{}, nil
}

func (s *productService) GetImages(productID int) ([]models.ProductImage, error) {
	galleryImages, err := s.productRepository.GetGalleryImages(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery images: %w", err)
	}

	var images []models.ProductImage
	for _, galleryImage := range galleryImages {
		title, err := s.productRepository.ProductFormName(productID)
		if err != nil {
			//log.Printf("Failed to get product form name for product %d: %v", productID, err)
			continue
		}

		ext := "png"
		if galleryImage.Webp.Bool {
			ext = "webp"
		}
		image := models.ProductImage{}
		image.Description = galleryImage.Description.String
		image.Links.Medium = fmt.Sprintf("/images/product/gallery/%d/%d/medium.%s", galleryImage.OwnerId, galleryImage.Id, ext)
		image.Links.Thumbnail = fmt.Sprintf("/images/product/gallery/%d/%d/preview.%s", galleryImage.OwnerId, galleryImage.Id, ext)
		image.Links.Raw = fmt.Sprintf("/images/product/gallery/%d/%d/original.%s", galleryImage.OwnerId, galleryImage.Id, ext)
		image.Links.OgImage = fmt.Sprintf("/images/product/gallery/%d/%d/og_image.webp", galleryImage.OwnerId, galleryImage.Id)

		if galleryImage.Name.Valid {
			image.Title = fmt.Sprintf("%s %s в Казахстане, интернет-аптека Рокет Фарм", galleryImage.Name.String, title)
		}

		images = append(images, image)
	}

	return images, nil
}

func (s *productService) SetAllProductToElastic() error {
	offset := 0
	limit := 10000
	for {
		products, err := s.productRepository.AllProductSearchData(offset, limit)
		if err != nil {
			return fmt.Errorf("failed to get products: %w", err)
		}
		if len(products) == 0 {
			break
		}

		esProducts := make([]models.ProductElastic, 0, len(products))
		for _, product := range products {
			log.Printf("Indexing product: %d", product.Id)
			productPharmacy, _ := s.productRepository.ProductPharmacy(product.Id)
			esProducts = append(esProducts, product.ToProductElastic(productPharmacy.Price, productPharmacy.Count))
		}

		if err := s.elastic.ProductAddDocument(esProducts); err != nil {
			log.Printf("Error indexing products to Elasticsearch: %v", err)
			// Consider whether to return error here or continue with next batch
		}

		offset += limit
	}

	return nil
}
