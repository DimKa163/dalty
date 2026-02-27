package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltymodel"
)

type (
	Spec struct {
		ID            string
		IntegrationID string `json:"integration_id"`
		Fnrec         string `json:"fnrec"`
		Quantity      int32  `json:"quantity"`
		RelateToID    string `json:"relate_to_id"`
	}
	SpecRequest struct {
		Specs map[string]*Spec `json:"specs"`
	}
	SpecResponse struct {
		Specifications []*daltymodel.Specification `json:"specifications"`
	}
	ProductSpec struct {
		ID string `json:"id"`
		*core.Product
		Quantity   int32  `json:"quantity"`
		RelateToID string `json:"relate_to_id"`
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
	reverseSpecs := make(map[string]*ProductSpec)
	for _, spec := range prdSpecs {
		rel, err := ss.relationRepository.GetByLeftID(ctx, spec.Product.ID)
		if err != nil {
			return nil, err
		}
		if len(rel) == 0 {
			reverseSpecs[spec.ID] = spec
			continue
		}
		// TODO evaluate direct relation
		subProducts := make([]*ProductSpec, len(rel))
		for i, product := range rel {
			subProducts[i] = &ProductSpec{ID: "", Product: product.Right, Quantity: spec.Quantity * product.Amount}
		}
		specs = append(
			specs,
			toDirectDaltySpecification(spec, subProducts),
		)
	}
	defaultSpecs := make([]*ProductSpec, 0, len(prdSpecs)/2)
	for _, left := range reverseSpecs {
		if left.RelateToID == "" {
			continue
		}
		if left.Quantity == 0 {
			continue
		}
		right, ok := reverseSpecs[left.RelateToID]
		if !ok {
			return nil, fmt.Errorf("no specification found for product %q", left.ID)
		}
		specs = append(specs, toReverseSpecification([]*ProductSpec{left, right}))
		left.Quantity = 0
		right.Quantity = 0
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
			daltymodel.NewLine(spec.ID, toDaltyProduct(spec.Product), spec.Quantity, daltymodel.PickupStrategyNearest),
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

func (ss *SpecificationService) getProductSpec(ctx context.Context, specs map[string]*Spec) (map[string]*ProductSpec, error) {
	products := make(map[string]*ProductSpec)
	errs := make([]*daltyerrors.StorageError, 0)
	for i, spec := range specs {
		prd, err := ss.find(ctx, spec)
		if err != nil {

			var storageErr *daltyerrors.StorageError
			if errors.As(err, &storageErr) {
				errs = append(errs, storageErr)
			}
		}
		products[i] = &ProductSpec{ID: i, Product: prd, Quantity: spec.Quantity, RelateToID: spec.RelateToID}
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

func toDirectDaltySpecification(head *ProductSpec, products []*ProductSpec) *daltymodel.Specification {
	h := toDaltyProduct(head.Product)
	children := make([]*daltymodel.Line, len(products))
	for i, product := range products {
		children[i] = toDaltyLine(product.ID, toDaltyProduct(product.Product), product.Quantity, func(prd *daltymodel.Product) daltymodel.PickupStrategy {
			switch product.Group {
			case daltymodel.ProductGroupChildrenBedBases:
				return daltymodel.PickupStrategyFarthest
			case daltymodel.ProductGroupSlattedBases:
				return daltymodel.PickupStrategyFarthest
			case daltymodel.ProductGroupBedBasesWithStorage:
				return daltymodel.PickupStrategyFarthest
			default:
				return daltymodel.PickupStrategyNearest
			}
		})
	}
	return daltymodel.NewDirectSpecification(daltymodel.NewLine(
		head.ID,
		h,
		head.Quantity,
		daltymodel.PickupStrategyNearest,
	), daltymodel.PickupStrategyNearest, children)
}

func toReverseSpecification(products []*ProductSpec) *daltymodel.Specification {
	children := make([]*daltymodel.Line, len(products))
	for i, product := range products {
		children[i] = toDaltyLine(product.ID, toDaltyProduct(product.Product), product.Quantity, func(prd *daltymodel.Product) daltymodel.PickupStrategy {
			switch product.Group {
			case daltymodel.ProductGroupChildrenBedBases:
				return daltymodel.PickupStrategyFarthest
			case daltymodel.ProductGroupSlattedBases:
				return daltymodel.PickupStrategyFarthest
			case daltymodel.ProductGroupBedBasesWithStorage:
				return daltymodel.PickupStrategyFarthest
			default:
				return daltymodel.PickupStrategyNearest
			}
		})
	}
	return daltymodel.NewReverseSpecification(nil, daltymodel.PickupStrategyNearest, children)
}

func toDaltyLine(id string, product *daltymodel.Product, quantity int32, fn func(prd *daltymodel.Product) daltymodel.PickupStrategy) *daltymodel.Line {
	return daltymodel.NewLine(id, product, quantity, fn(product))
}

func toDaltyProduct(product *core.Product) *daltymodel.Product {
	v := daltymodel.Product(*product)
	return &v
}

func validateSpecs(products map[string]*ProductSpec) error {
	for _, product := range products {
		if product.IsArchive {
			return daltyerrors.New(5, &daltyerrors.EntityError{ID: product.Product.ID.String(), EntityName: "product"})
		}
	}
	return nil
}
