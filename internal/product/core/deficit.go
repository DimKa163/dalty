package core

import (
	"context"
	"time"

	"github.com/beevik/guid"
)

type Deficit struct {
	ID           guid.Guid
	ProductID    guid.Guid
	PurchaseDate *time.Time
	WarehouseID  *guid.Guid
}

type DeficitRepository interface {
	GetByProductID(ctx context.Context, productID guid.Guid) (*Deficit, error)
}
