package server

import (
	"context"
	pb "github.com/antibomberman/aurma-protos/gen/go/product"
)

func (s server) Search(ctx context.Context, req *pb.ProductSearchRequest) (*pb.ProductSearchResponse, error) {
	productDetails, total, err := s.productService.Search(req.Title, int(req.Page), int(req.Limit), req.Sort.String(), int(req.MinPrice), int(req.MaxPrice))
	if err != nil {
		return nil, err
	}
	products := make([]*pb.Product, len(productDetails))
	for i, pd := range productDetails {
		products[i] = pd.ToPbProduct()
	}
	return &pb.ProductSearchResponse{
		Products:   products,
		TotalCount: int32(total),
	}, nil
}
