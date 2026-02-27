package usecase

import (
	"context"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/internal/product/mocks"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"github.com/beevik/guid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecuteDirectSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockProductRepository := mocks.NewMockProductRepository(ctrl)
	mockRelateRepository := mocks.NewMockRelationRepository(ctrl)
	product := &core.Product{
		ID:             *guid.New(),
		Name:           "Test",
		Group:          daltymodel.ProductGroupArmchairs,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "1",
		Fnrec:          "1",
	}
	relation := map[guid.Guid][]*core.Relation{
		product.ID: {
			&core.Relation{
				ID:      *guid.New(),
				RightID: *guid.New(),
				Right: &core.Product{
					ID:            *guid.New(),
					Name:          "Sub product 1",
					Group:         daltymodel.ProductGroupArmchairs,
					IntegrationID: "11",
					Fnrec:         "11",
				},
			},
			&core.Relation{
				ID:      *guid.New(),
				RightID: *guid.New(),
				Right: &core.Product{
					ID:            *guid.New(),
					Name:          "Sub product 2",
					Group:         daltymodel.ProductGroupArmchairs,
					IntegrationID: "12",
					Fnrec:         "12",
				},
			},
		},
	}
	req := &SpecRequest{
		Specs: make(map[string]*Spec),
	}
	specID := guid.NewString()
	req.Specs[specID] = &Spec{
		ID:            specID,
		IntegrationID: product.IntegrationID,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, product.IntegrationID).Return(product, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	specs, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, specs, "")

	for _, spec := range specs {
		assert.Equal(t, daltymodel.PickupStrategyNearest, spec.Strategy)
		assert.Equal(t, daltymodel.SpecificationTypeDirect, spec.Type)
		for _, l := range spec.ChildProducts {
			switch l.Product.Group {
			case daltymodel.ProductGroupChildrenBedBases:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			case daltymodel.ProductGroupSlattedBases:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			case daltymodel.ProductGroupBedBasesWithStorage:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			default:
				assert.Equal(t, daltymodel.PickupStrategyNearest, l.Strategy)
			}
		}
	}
}

func TestExecuteReverseSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockProductRepository := mocks.NewMockProductRepository(ctrl)
	mockRelateRepository := mocks.NewMockRelationRepository(ctrl)
	productA := &core.Product{
		ID:             *guid.New(),
		Name:           "Test A",
		Group:          daltymodel.ProductGroupBeds,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "1",
		Fnrec:          "1",
	}
	productB := &core.Product{
		ID:             *guid.New(),
		Name:           "Test B",
		Group:          daltymodel.ProductGroupBedBases,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "2",
		Fnrec:          "2",
	}
	relation := map[guid.Guid][]*core.Relation{
		productA.ID: make([]*core.Relation, 0),
		productB.ID: make([]*core.Relation, 0),
	}
	req := &SpecRequest{
		Specs: make(map[string]*Spec),
	}
	specID1 := guid.NewString()
	specID2 := guid.NewString()
	req.Specs[specID1] = &Spec{
		ID:            specID1,
		IntegrationID: productA.IntegrationID,
		RelateToID:    specID2,
		Quantity:      1,
	}
	req.Specs[specID2] = &Spec{
		ID:            specID2,
		IntegrationID: productB.IntegrationID,
		RelateToID:    specID1,
		Quantity:      1,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productA.IntegrationID).Return(productA, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productB.IntegrationID).Return(productB, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	specs, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, specs, "")
	assert.Equal(t, 1, len(specs))
	for _, spec := range specs {
		assert.Equal(t, daltymodel.PickupStrategyNearest, spec.Strategy)
		assert.Equal(t, daltymodel.SpecificationTypeReverse, spec.Type)
		for _, l := range spec.ChildProducts {
			switch l.Product.Group {
			case daltymodel.ProductGroupChildrenBedBases:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			case daltymodel.ProductGroupSlattedBases:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			case daltymodel.ProductGroupBedBasesWithStorage:
				assert.Equal(t, daltymodel.PickupStrategyFarthest, l.Strategy)
			default:
				assert.Equal(t, daltymodel.PickupStrategyNearest, l.Strategy)
			}
		}
	}
}

func TestExecuteDefaultSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockProductRepository := mocks.NewMockProductRepository(ctrl)
	mockRelateRepository := mocks.NewMockRelationRepository(ctrl)
	productA := &core.Product{
		ID:             *guid.New(),
		Name:           "Test A",
		Group:          daltymodel.ProductGroupBeds,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "1",
		Fnrec:          "1",
	}
	productB := &core.Product{
		ID:             *guid.New(),
		Name:           "Test B",
		Group:          daltymodel.ProductGroupBedBases,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "2",
		Fnrec:          "2",
	}
	relation := map[guid.Guid][]*core.Relation{
		productA.ID: make([]*core.Relation, 0),
		productB.ID: make([]*core.Relation, 0),
	}
	req := &SpecRequest{
		Specs: make(map[string]*Spec, 2),
	}
	specID1 := guid.NewString()
	specID2 := guid.NewString()

	req.Specs[specID1] = &Spec{
		IntegrationID: productA.IntegrationID,
		Quantity:      1,
	}
	req.Specs[specID2] = &Spec{
		IntegrationID: productB.IntegrationID,
		Quantity:      1,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productA.IntegrationID).Return(productA, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productB.IntegrationID).Return(productB, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	specs, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, specs, "")
	assert.Equal(t, 2, len(specs))

	for _, spec := range specs {
		assert.Equal(t, daltymodel.PickupStrategyNearest, spec.Strategy)
		assert.Equal(t, daltymodel.SpecificationTypeDefault, spec.Type)
		assert.Equal(t, daltymodel.PickupStrategyNearest, spec.Product.Strategy)
		assert.Empty(t, spec.ChildProducts)
	}
}

func TestExecuteCombineSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockProductRepository := mocks.NewMockProductRepository(ctrl)
	mockRelateRepository := mocks.NewMockRelationRepository(ctrl)
	productA := &core.Product{
		ID:             *guid.New(),
		Name:           "Test A",
		Group:          daltymodel.ProductGroupBeds,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "1",
		Fnrec:          "1",
	}
	productB := &core.Product{
		ID:             *guid.New(),
		Name:           "Test B",
		Group:          daltymodel.ProductGroupBeds,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "2",
		Fnrec:          "2",
	}
	productC := &core.Product{
		ID:             *guid.New(),
		Name:           "Test C",
		Group:          daltymodel.ProductGroupBedBases,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "3",
		Fnrec:          "3",
	}
	productD := &core.Product{
		ID:             *guid.New(),
		Name:           "Test D",
		Group:          daltymodel.ProductGroupBeds,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "4",
		Fnrec:          "4",
	}
	productE := &core.Product{
		ID:             *guid.New(),
		Name:           "Test E",
		Group:          daltymodel.ProductGroupSlattedBases,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "5",
		Fnrec:          "5",
	}
	productF := &core.Product{
		ID:             *guid.New(),
		Name:           "Test F",
		Group:          daltymodel.ProductGroupArmchairs,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "6",
		Fnrec:          "6",
	}
	productG := &core.Product{
		ID:             *guid.New(),
		Name:           "Test G",
		Group:          daltymodel.ProductGroupBedBases,
		ProductionType: core.ProductionTypeProducing,
		IntegrationID:  "7",
		Fnrec:          "7",
	}
	relation := map[guid.Guid][]*core.Relation{
		productA.ID: {
			&core.Relation{
				ID:      *guid.New(),
				RightID: *guid.New(),
				Right: &core.Product{
					ID:            *guid.New(),
					Name:          "Sub product A1",
					Group:         daltymodel.ProductGroupBeds,
					IntegrationID: "11",
					Fnrec:         "11",
				},
				Amount: 1,
			},
			&core.Relation{
				ID:      *guid.New(),
				RightID: *guid.New(),
				Right: &core.Product{
					ID:            *guid.New(),
					Name:          "Sub product A2",
					Group:         daltymodel.ProductGroupBedBases,
					IntegrationID: "12",
					Fnrec:         "12",
				},
				Amount: 2,
			},
		},
		productB.ID: make([]*core.Relation, 0),
		productC.ID: make([]*core.Relation, 0),
		productD.ID: make([]*core.Relation, 0),
		productE.ID: make([]*core.Relation, 0),
		productF.ID: make([]*core.Relation, 0),
		productG.ID: make([]*core.Relation, 0),
	}
	req := &SpecRequest{
		Specs: make(map[string]*Spec),
	}
	specID1 := guid.NewString()
	specID2 := guid.NewString()
	specID3 := guid.NewString()
	specID4 := guid.NewString()
	specID5 := guid.NewString()
	specID6 := guid.NewString()
	specID7 := guid.NewString()
	req.Specs[specID1] = &Spec{
		IntegrationID: productA.IntegrationID,
		Quantity:      1,
	}
	req.Specs[specID2] = &Spec{
		IntegrationID: productB.IntegrationID,
		Quantity:      1,
		RelateToID:    specID3,
	}
	req.Specs[specID3] = &Spec{
		IntegrationID: productC.IntegrationID,
		Quantity:      1,
		RelateToID:    specID2,
	}
	req.Specs[specID4] = &Spec{
		IntegrationID: productD.IntegrationID,
		Quantity:      1,
		RelateToID:    specID5,
	}
	req.Specs[specID5] = &Spec{
		IntegrationID: productE.IntegrationID,
		Quantity:      3,
		RelateToID:    specID4,
	}
	req.Specs[specID6] = &Spec{
		IntegrationID: productF.IntegrationID,
		Quantity:      1,
	}
	req.Specs[specID7] = &Spec{
		IntegrationID: productG.IntegrationID,
		Quantity:      1,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productA.IntegrationID).Return(productA, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productB.IntegrationID).Return(productB, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productC.IntegrationID).Return(productC, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productD.IntegrationID).Return(productD, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productE.IntegrationID).Return(productE, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productF.IntegrationID).Return(productF, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productG.IntegrationID).Return(productG, nil)
	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	specs, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, specs, "")
	assert.Equal(t, 5, len(specs))

	assert.Equal(t, daltymodel.SpecificationTypeDirect, specs[0].Type)
	assert.Equal(t, 2, len(specs[0].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeReverse, specs[1].Type)
	assert.Equal(t, 2, len(specs[1].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeReverse, specs[2].Type)
	assert.Equal(t, 2, len(specs[2].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeDefault, specs[3].Type)
	assert.Equal(t, 0, len(specs[3].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeDefault, specs[4].Type)
	assert.Equal(t, 0, len(specs[4].ChildProducts))
}
