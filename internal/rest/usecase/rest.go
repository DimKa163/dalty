package usecase

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/internal/rest/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/beevik/guid"
)

type (
	RestRequest struct {
		FilialID      guid.Guid   `json:"filial_id"`
		WarehousesIDS []guid.Guid `json:"warehouses_ids"`
		ProductIDS    []guid.Guid `json:"product_ids"`
	}
	RestResult struct {
		ProductMap map[guid.Guid]map[guid.Guid]*core.Rest
	}
)

type RestService struct {
	restRepository core.RestRepository
}

func NewRestService(restRepository core.RestRepository) *RestService {
	return &RestService{
		restRepository: restRepository,
	}
}

func (service *RestService) GetRest(ctx context.Context, request *RestRequest) (*RestResult, error) {
	var result RestResult
	result.ProductMap = make(map[guid.Guid]map[guid.Guid]*core.Rest)
	for _, warehouseID := range request.WarehousesIDS {
		rests := make(map[guid.Guid]*core.Rest)
		for _, productID := range request.ProductIDS {
			rest, err := service.restRepository.FindRest(ctx, productID, warehouseID, request.FilialID)
			if err != nil {
				if errors.Is(err, daltyerrors.ErrNotFound) {
					continue
				}
				return nil, err
			}
			rests[productID] = rest
		}
		result.ProductMap[warehouseID] = rests
	}
	return &result, nil
}
