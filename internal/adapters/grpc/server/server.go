package server

import (
	"aurma_product/internal/config"
	"aurma_product/internal/di"
	"aurma_product/internal/services"
	pb "github.com/antibomberman/aurma-protos/gen/go/product"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
	productService services.ProductService
	cfg            *config.Config
}

func Register(gRPC *grpc.Server, container *di.Container) {
	pb.RegisterProductServiceServer(gRPC, &server{
		productService: container.ProductService,
		cfg:            container.Config,
	})
}
