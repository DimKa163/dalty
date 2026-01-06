package usecase

import (
	"context"

	"github.com/DimKa163/dalty/internal/product/core"
)

type ProductRequest struct {
	IntegrationID string
	Fnrec         string
}
type ProductService struct {
	productRepository core.ProductRepository
}

func NewProductService(productRepository core.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (ps *ProductService) BatchRequest(ctx context.Context, requests []*ProductRequest) ([]*core.Product, error) {
	results := make([]*core.Product, 0, len(requests))
	for _, request := range requests {
		resp, err := ps.Find(ctx, request)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}
	return results, nil
}

func (ps *ProductService) Find(ctx context.Context, request *ProductRequest) (*core.Product, error) {
	if request.IntegrationID != "" {
		product, err := ps.productRepository.GetByIntegrationID(ctx, request.IntegrationID)
		if err != nil {
			return nil, err
		}
		return product, nil
	} else {
		product, err := ps.productRepository.GetByFnrec(ctx, request.Fnrec)
		if err != nil {
			return nil, err
		}
		return product, nil
	}
}
