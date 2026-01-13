package usecase

import (
	"context"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/internal/product/mocks"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
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
		Specs: make([]*Spec, 1),
	}
	req.Specs[0] = &Spec{
		IntegrationID: product.IntegrationID,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, product.IntegrationID).Return(product, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	r, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, r, "")
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
		Specs: make([]*Spec, 2),
	}
	req.Specs[0] = &Spec{
		IntegrationID: productA.IntegrationID,
		Quantity:      1,
	}
	req.Specs[1] = &Spec{
		IntegrationID: productB.IntegrationID,
		Quantity:      1,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productA.IntegrationID).Return(productA, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productB.IntegrationID).Return(productB, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}
	left := *guid.New()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productA.ID, productB.ID).Return(&core.Relation{
		ID:      *guid.New(),
		LeftID:  left,
		RightID: productA.ID,
		Amount:  1,
	}, &core.Relation{
		ID:      *guid.New(),
		LeftID:  left,
		RightID: productB.ID,
		Amount:  1,
	}, nil)

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	r, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, r, "")
	assert.Equal(t, 1, len(r))
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
		Specs: make([]*Spec, 2),
	}
	req.Specs[0] = &Spec{
		IntegrationID: productA.IntegrationID,
		Quantity:      1,
	}
	req.Specs[1] = &Spec{
		IntegrationID: productB.IntegrationID,
		Quantity:      1,
	}
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productA.IntegrationID).Return(productA, nil)
	mockProductRepository.EXPECT().GetByIntegrationID(ctx, productB.IntegrationID).Return(productB, nil)

	for i, r := range relation {
		mockRelateRepository.EXPECT().GetByLeftID(ctx, i).Return(r, nil)
	}
	mockRelateRepository.EXPECT().GetByRightID(ctx, productA.ID, productB.ID).Return(nil, nil, daltyerrors.ErrNotFound).Times(1)

	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	r, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, r, "")
	assert.Equal(t, 2, len(r))
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
		Specs: make([]*Spec, 7),
	}
	req.Specs[0] = &Spec{
		IntegrationID: productA.IntegrationID,
		Quantity:      1,
	}
	req.Specs[1] = &Spec{
		IntegrationID: productB.IntegrationID,
		Quantity:      1,
	}
	req.Specs[2] = &Spec{
		IntegrationID: productC.IntegrationID,
		Quantity:      1,
	}
	req.Specs[3] = &Spec{
		IntegrationID: productD.IntegrationID,
		Quantity:      1,
	}
	req.Specs[4] = &Spec{
		IntegrationID: productE.IntegrationID,
		Quantity:      3,
	}
	req.Specs[5] = &Spec{
		IntegrationID: productF.IntegrationID,
		Quantity:      1,
	}
	req.Specs[6] = &Spec{
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
	left := &core.Product{
		ID:   *guid.New(),
		Name: "Head BC",
	}
	mockRelateRepository.EXPECT().GetByRightID(ctx, productB.ID, productC.ID).Return(&core.Relation{
		ID:      *guid.New(),
		LeftID:  left.ID,
		Left:    left,
		RightID: productB.ID,
		Amount:  1,
	}, &core.Relation{
		ID:      *guid.New(),
		LeftID:  left.ID,
		Left:    left,
		RightID: productC.ID,
		Amount:  1,
	}, nil)
	mockRelateRepository.EXPECT().GetByRightID(ctx, productB.ID, productD.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productB.ID, productE.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productB.ID, productF.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productB.ID, productG.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()

	mockRelateRepository.EXPECT().GetByRightID(ctx, productC.ID, productD.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productC.ID, productE.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productC.ID, productF.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productC.ID, productG.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()

	left = &core.Product{
		ID:   *guid.New(),
		Name: "Head DE",
	}
	mockRelateRepository.EXPECT().GetByRightID(ctx, productD.ID, productE.ID).Return(&core.Relation{
		ID:      *guid.New(),
		LeftID:  left.ID,
		Left:    left,
		RightID: productD.ID,
		Amount:  1,
	}, &core.Relation{
		ID:      *guid.New(),
		LeftID:  left.ID,
		Left:    left,
		RightID: productE.ID,
		Amount:  2,
	}, nil)
	mockRelateRepository.EXPECT().GetByRightID(ctx, productD.ID, productF.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productD.ID, productG.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()

	mockRelateRepository.EXPECT().GetByRightID(ctx, productE.ID, productF.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	mockRelateRepository.EXPECT().GetByRightID(ctx, productE.ID, productG.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()

	mockRelateRepository.EXPECT().GetByRightID(ctx, productF.ID, productG.ID).Return(nil, nil, daltyerrors.ErrNotFound).AnyTimes()
	sut := NewSpecificationService(mockProductRepository, mockRelateRepository)

	r, err := sut.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, r, "")
	assert.Equal(t, 6, len(r))

	assert.Equal(t, daltymodel.SpecificationTypeDirect, r[0].Type)
	assert.Equal(t, 2, len(r[0].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeReverse, r[1].Type)
	assert.Equal(t, 2, len(r[1].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeReverse, r[2].Type)
	assert.Equal(t, 2, len(r[2].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeDefault, r[3].Type)
	assert.Equal(t, 0, len(r[3].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeDefault, r[4].Type)
	assert.Equal(t, 0, len(r[4].ChildProducts))
	assert.Equal(t, daltymodel.SpecificationTypeDefault, r[5].Type)
	assert.Equal(t, 0, len(r[5].ChildProducts))
}
