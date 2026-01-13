package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"math"
)

type (
	Spec struct {
		IntegrationID string `json:"integration_id"`
		Fnrec         string `json:"fnrec"`
		Quantity      int32  `json:"quantity"`
	}
	SpecRequest struct {
		Specs []*Spec `json:"specs"`
	}
	SpecResponse struct {
		Specifications []*daltymodel.Specification `json:"specifications"`
	}
	ProductSpec struct {
		*core.Product
		Quantity int32 `json:"quantity"`
	}
)
type SpecificationService struct {
	productRepository  core.ProductRepository
	relationRepository core.RelationRepository
}

func NewSpecificationService(productRepository core.ProductRepository, relationRepository core.RelationRepository) *SpecificationService {
	return &SpecificationService{productRepository, relationRepository}
}

func (ss *SpecificationService) Execute(ctx context.Context, request *SpecRequest) ([]*daltymodel.Specification, error) {
	prdSpecs, err := ss.getProductSpec(ctx, request.Specs)
	if err != nil {
		return nil, err
	}
	specs := make([]*daltymodel.Specification, 0, len(prdSpecs))
	reverseSpecs := make([]*ProductSpec, 0, len(prdSpecs))
	for _, spec := range prdSpecs {
		rel, err := ss.relationRepository.GetByLeftID(ctx, spec.ID)
		if err != nil {
			return nil, err
		}
		if len(rel) == 0 {
			reverseSpecs = append(reverseSpecs, spec)
			continue
		}
		// TODO evaluate direct relation
		subProducts := make([]*daltymodel.Line, len(rel))
		for i, product := range rel {
			subProducts[i] = daltymodel.NewLine(toDaltyProduct(product.Right), product.Amount*spec.Quantity, daltymodel.PickupStrategyNearest)
		}
		specs = append(
			specs,
			daltymodel.NewDirectSpecification(daltymodel.NewLine(toDaltyProduct(spec.Product),
				spec.Quantity,
				daltymodel.PickupStrategyNearest),
				daltymodel.PickupStrategyNearest,
				subProducts),
		)
	}
	defaultSpecs := make([]*ProductSpec, 0, len(prdSpecs)/2)
	for i, left := range reverseSpecs {
		if left.Quantity == 0 {
			continue
		}
		for j := i + 1; j < len(reverseSpecs); j++ {
			right := reverseSpecs[j]
			if left.ID == right.ID {
				continue
			}
			if right.Quantity == 0 {
				continue
			}
			l, r, err := ss.relationRepository.GetByRightID(ctx, left.ID, right.ID)
			if err != nil {
				if errors.Is(err, daltyerrors.ErrNotFound) {
					continue
				}
				return nil, err
			}

			lower := math.Min(float64(left.Quantity), float64(right.Quantity))
			for lower > 0 {

				lq := left.Quantity - l.Amount
				rq := right.Quantity - r.Amount
				if lq < 0 || rq < 0 {
					break
				}
				// TODO evaluate reverse specification
				specs = append(specs,
					daltymodel.NewReverseSpecification(daltymodel.NewLine(toDaltyProduct(l.Left), 1, daltymodel.PickupStrategyNearest),
						daltymodel.PickupStrategyNearest, []*daltymodel.Line{
							daltymodel.NewLine(toDaltyProduct(left.Product), left.Quantity, daltymodel.PickupStrategyNearest),
							daltymodel.NewLine(toDaltyProduct(right.Product), right.Quantity, daltymodel.PickupStrategyNearest),
						}))
				left.Quantity = lq
				right.Quantity = rq
				lower = math.Min(float64(left.Quantity), float64(right.Quantity))
			}
		}
	}
	for _, spec := range reverseSpecs {
		if spec.Quantity == 0 {
			continue
		}
		defaultSpecs = append(defaultSpecs, spec)
	}
	for _, spec := range defaultSpecs {
		// TODO evaluate default
		specs = append(specs, daltymodel.NewDefaultSpecification(
			daltymodel.NewLine(toDaltyProduct(spec.Product), spec.Quantity, daltymodel.PickupStrategyNearest),
			daltymodel.PickupStrategyNearest))
	}
	return specs, nil
}

func (ss *SpecificationService) getByLeft(ctx context.Context, spec *Spec) ([]*core.Relation, error) {
	var fn func(context.Context, string) ([]*core.Relation, error)
	var filterVal string
	switch {
	case spec.IntegrationID != "":
		filterVal = spec.IntegrationID
		fn = ss.relationRepository.GetByLeftIntegrationID
	case spec.Fnrec != "":
		filterVal = spec.Fnrec
		fn = ss.relationRepository.GetByLeftFnrec
	default:
		return nil, fmt.Errorf("no ids provided")
	}
	return fn(ctx, filterVal)
}

func (ss *SpecificationService) getProductSpec(ctx context.Context, specs []*Spec) ([]*ProductSpec, error) {
	products := make([]*ProductSpec, len(specs))
	errs := make([]*daltyerrors.StorageError, 0)
	for i, spec := range specs {
		prd, err := ss.find(ctx, spec)
		if err != nil {

			var storageErr *daltyerrors.StorageError
			if errors.As(err, &storageErr) {
				errs = append(errs, storageErr)
			}
		}
		products[i] = &ProductSpec{Product: prd, Quantity: spec.Quantity}
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
	if err := validateSpecs(products); err != nil {
		return nil, err
	}
	return products, nil
}

func (ss *SpecificationService) find(ctx context.Context, request *Spec) (*core.Product, error) {
	var productFunc func(context.Context, string) (*core.Product, error)
	var filter string
	if request.IntegrationID != "" {
		productFunc = ss.productRepository.GetByIntegrationID
		filter = request.IntegrationID
	} else {
		productFunc = ss.productRepository.GetByFnrec
		filter = request.Fnrec
	}
	return productFunc(ctx, filter)
}

func toDaltyProduct(product *core.Product) *daltymodel.Product {
	v := daltymodel.Product(*product)
	return &v
}

func validateSpecs(products []*ProductSpec) error {
	for _, product := range products {
		if product.IsArchive {
			return daltyerrors.New(5, &daltyerrors.EntityError{ID: product.ID.String(), EntityName: "product"})
		}
	}
	return nil
}
