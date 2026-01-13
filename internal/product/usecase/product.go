package usecase

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
)

var ErrArchiveProduct = errors.New("product is archived")

var ErrProductsNotFound = errors.New("products not found")

type ProductRequest struct {
	ID            string
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
	errs := make([]*daltyerrors.StorageError, 0)
	for _, request := range requests {
		resp, err := ps.find(ctx, request)
		if err != nil {

			var storageErr *daltyerrors.StorageError
			if errors.As(err, &storageErr) {
				errs = append(errs, storageErr)
			}
		}
		results = append(results, resp)
	}
	if len(errs) > 0 {
		daltyErrs := make([]*daltyerrors.EntityError, len(errs))
		for i, err := range errs {
			daltyErrs[i] = &daltyerrors.EntityError{
				ID:         err.Value[0].(string),
				EntityName: "product",
			}
		}
		return nil, daltyerrors.New(
			6,
			daltyErrs...,
		)
	}
	if err := validateProducts(results); err != nil {
		return nil, err
	}
	return results, nil
}

func (ps *ProductService) Find(ctx context.Context, request *ProductRequest) (*core.Product, error) {
	product, err := ps.find(ctx, request)
	if err != nil {
		var storageErr *daltyerrors.StorageError
		if errors.As(err, &storageErr) {
			return nil, daltyerrors.New(
				6,
				&daltyerrors.EntityError{ID: storageErr.Value[0].(string), EntityName: "product"},
			)
		}
	}
	if err = validateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (ps *ProductService) find(ctx context.Context, request *ProductRequest) (*core.Product, error) {
	var productFunc func(context.Context, string) (*core.Product, error)
	var filter string
	if request.ID != "" {
		productFunc = ps.productRepository.GetByID
		filter = request.ID
	} else if request.IntegrationID != "" {
		productFunc = ps.productRepository.GetByIntegrationID
		filter = request.IntegrationID
	} else {
		productFunc = ps.productRepository.GetByFnrec
		filter = request.Fnrec
	}
	return productFunc(ctx, filter)
}

func validateProduct(product *core.Product) error {
	if product.IsArchive {
		return daltyerrors.New(5, &daltyerrors.EntityError{ID: product.ID.String(), EntityName: "product"})
	}
	return nil
}

func validateProducts(products []*core.Product) error {
	for _, product := range products {
		if product.IsArchive {
			return daltyerrors.New(5, &daltyerrors.EntityError{ID: product.ID.String(), EntityName: "product"})
		}
	}
	return nil
}
