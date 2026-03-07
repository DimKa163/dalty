package core

import (
	"context"
	"github.com/beevik/guid"
	"github.com/shopspring/decimal"
)

type Rest struct {
	ID             guid.Guid       `json:"id"`
	ProductID      guid.Guid       `json:"product_id"`
	SubWarehouseID guid.Guid       `json:"sub_warehouse_id"`
	FilialID       guid.Guid       `json:"filial_id"`
	IntegrationID  string          `json:"integration_id"`
	Quantity       decimal.Decimal `json:"quantity"`
}

type RestRepository interface {
	FindRest(ctx context.Context, productID, subWarehouse, filialID guid.Guid) (*Rest, error)
}
