package di

import (
	"aurma_product/internal/config"
	"aurma_product/internal/database"
	"aurma_product/internal/elastic"
	"aurma_product/internal/repositories"
	"aurma_product/internal/services"
	"github.com/antibomberman/dblayer"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	Config          *config.Config
	DB              *sqlx.DB
	Elastic         *elastic.Elastic
	ProductService  services.ProductService
	SadykhanService services.SadykhanService
}

func NewContainer() (*Container, error) {
	container := &Container{}
	var err error

	// Load config
	container.Config = config.Load()

	// Initialize DB
	container.DB, err = database.Connect(container.Config)
	if err != nil {
		return nil, err
	}
	// Initialize Elasticsearch
	container.Elastic, err = elastic.New()
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	productRepo := repositories.NewProductRepository(container.DB)
	dblayer := dblayer.NewDBLayer(container.DB)

	// Initialize services
	container.ProductService = services.NewProductService(dblayer, productRepo, container.Elastic)
	container.SadykhanService = services.NewSadykhanService(dblayer, productRepo)

	return container, nil
}

func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
