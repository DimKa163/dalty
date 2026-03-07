package daltymodel

import "github.com/beevik/guid"

type WarehouseType int

const (
	NodeUnrecognized WarehouseType = iota
	NodeFree
	NodeMain
	NodeCenter
	NodeMall
	NodeTransit
	NodeReservation
	NodeLoses
	NodeMarketing
	NodeExposition
	NodePartner
	NodePartner2
	NodeFree2
	NodeProblem
	NodeRefund
	NodeProduction
	NodeRecycling
	NodeService
	NodeMaterial
	NodeMarkdown
	NodeBuffer
	NodeDiscount
	NodeCentralMainIntermediate
	NodeMainCentralIntermediate
	NodeCentralFreeIntermediate
	NodeFreeCentralIntermediate
)

func (w WarehouseType) String() string {
	names := []string{
		"UNRECOGNIZED",
		"FREE",
		"MAIN",
		"CENTER",
		"MALL",
		"TRANSIT",
		"RESERVATION",
		"LOSES",
		"MARKETING",
		"EXPOSITION",
		"PARTNER",
		"PARTNER2",
		"FREE2",
		"PROBLEM",
		"REFUND",
		"PRODUCTION",
		"RECYCLING",
		"SERVICE",
		"MATERIAL",
		"MARKDOWN",
		"BUFFER",
		"DISCOUNT",
		"CENTRAL_MAIN_INTERMEDIATE",
		"MAIN_CENTRAL_INTERMEDIATE",
		"CENTRAL_FREE_INTERMEDIATE",
		"FREE_CENTRAL_INTERMEDIATE",
	}

	if int(w) < 0 || int(w) >= len(names) {
		return "UNRECOGNIZED"
	}
	return names[w]
}

type Warehouse struct {
	ID                     guid.Guid     `json:"id"`
	Name                   string        `json:"name"`
	WarehouseType          WarehouseType `json:"warehouse_type"`
	TimeZone               string        `json:"time_zone"`
	AvailableRest          bool          `json:"available_rest"`
	Level                  int           `json:"level"`
	Address                string        `json:"address"`
	OnlyStockPickupAllowed bool          `json:"only_stock_pickup_allowed"`
	DescriptionGroup       string        `json:"description_group"`
}

type Path struct {
	Warehouses []*Warehouse `json:"warehouses"`
}

func (p *Path) Range(direct PickupStrategy, f func(w *Warehouse) error) error {
	switch direct {
	case PickupStrategyNearest:
		for i := 0; i < len(p.Warehouses); i++ {
			if err := f(p.Warehouses[i]); err != nil {
				return err
			}
		}
	case PickupStrategyFarthest:
		for i := len(p.Warehouses) - 1; i >= 0; i-- {
			if err := f(p.Warehouses[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
