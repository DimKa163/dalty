package core

import (
	"context"
	"database/sql"
	"github.com/DimKa163/dalty/internal/shared"

	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
)

type Warehouse struct {
	ID                     guid.Guid
	Fnrec                  string
	Name                   string
	IsActive               bool
	OnlyStockPickupAllowed bool
	SenderID               *guid.Guid
	RecipientID            *guid.Guid
	Type                   WarehouseType
	AvailableForBalance    bool
	Info                   *WarehouseInfo
}

func (w *Warehouse) Scan(dest pgx.Rows) error {
	var warehouseID string
	var warehouseFnrec string
	var name string
	var isActive bool
	var onlyStockPickupAllowed bool
	var senderID sql.NullString
	var recipientID sql.NullString
	var categoryFnrec sql.NullString
	var availableForBalances bool
	var warehouseInfoID sql.NullString
	var warehouseInfoFnrec sql.NullString
	var warehouseInfoAddress sql.NullString
	var warehouseInfoDescriptorGroup sql.NullString
	var tzID sql.NullString
	var tzCode sql.NullString
	if err := dest.Scan(&warehouseID,
		&warehouseFnrec,
		&name,
		&isActive,
		&onlyStockPickupAllowed,
		&senderID,
		&recipientID,
		&categoryFnrec,
		&availableForBalances,
		&warehouseInfoID,
		&warehouseInfoFnrec,
		&warehouseInfoAddress,
		&warehouseInfoDescriptorGroup,
		&tzID,
		&tzCode); err != nil {
		return err
	}
	var err error
	w.Name = name
	w.IsActive = isActive
	w.OnlyStockPickupAllowed = onlyStockPickupAllowed
	w.Type = MapWarehouseType(categoryFnrec.String)
	w.AvailableForBalance = availableForBalances
	var warehouseInfo *WarehouseInfo

	id, err := guid.ParseString(warehouseID)
	if err != nil {
		return err
	}
	w.ID = *id
	if senderID.Valid {
		w.SenderID, err = guid.ParseString(senderID.String)
		if err != nil {
			return err
		}
	}
	if recipientID.Valid {
		w.RecipientID, err = guid.ParseString(recipientID.String)
		if err != nil {
			return err
		}
	}
	if warehouseInfoID.Valid {
		warehouseInfo = &WarehouseInfo{}
		warehouseInfo.ID, err = guid.ParseString(warehouseInfoID.String)
		if err != nil {
			return err
		}
		var tz *TimeZone
		if tzID.Valid {
			tz = &TimeZone{}
			tz.ID, err = guid.ParseString(tzID.String)
			if err != nil {
				return err
			}
			tz.Code = tzCode.String
		}
		warehouseInfo.Address = warehouseInfoAddress.String
		warehouseInfo.Fnrec = warehouseInfoFnrec.String
		warehouseInfo.DescriptorGroup = warehouseInfoDescriptorGroup.String
		warehouseInfo.TimeZone = tz
	}
	w.Info = warehouseInfo
	return nil
}

type WarehouseInfo struct {
	ID              *guid.Guid
	Fnrec           string
	Address         string
	DescriptorGroup string
	TimeZone        *TimeZone
}

type TimeZone struct {
	ID   *guid.Guid
	Code string
}

type WarehouseRepository interface {
	GetAll(ctx context.Context) ([]*Warehouse, error)
}

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

func MapWarehouseType(code string) WarehouseType {
	switch code {
	case shared.WarehouseCategoryFree:
		return NodeFree
	case shared.WarehouseCategoryMain:
		return NodeMain
	case shared.WarehouseCategoryCentral:
		return NodeCenter
	case shared.WarehouseCategoryMall:
		return NodeMall
	case shared.WarehouseCategoryTransit:
		return NodeTransit
	case shared.WarehouseCategoryReservation:
		return NodeReservation
	case shared.WarehouseCategoryLoses:
		return NodeLoses
	case shared.WarehouseCategoryMarketing:
		return NodeMarketing
	case shared.WarehouseCategoryExposition:
		return NodeExposition
	case shared.WarehouseCategoryPartner:
		return NodePartner
	case shared.WarehouseCategoryPartner2:
		return NodePartner2
	case shared.WarehouseCategoryFree2:
		return NodeFree2
	case shared.WarehouseCategoryProduction:
		return NodeProduction
	case shared.WarehouseCategoryRecycling:
		return NodeRecycling
	case shared.WarehouseCategoryService:
		return NodeService
	case shared.WarehouseCategoryProblem:
		return NodeProblem
	case shared.WarehouseCategoryRefund:
		return NodeRefund
	case shared.WarehouseCategoryMaterial:
		return NodeMaterial
	case shared.WarehouseCategoryMarkdown:
		return NodeMarkdown
	case shared.WarehouseCategoryBuffer:
		return NodeBuffer
	case shared.WarehouseCategoryDiscount:
		return NodeDiscount
	case shared.WarehouseCategoryCentralMainIntermediate:
		return NodeCentralMainIntermediate
	case shared.WarehouseCategoryMainCentraIntermediate:
		return NodeMainCentralIntermediate
	case shared.WarehouseCategoryCentraFreeIntermediate:
		return NodeCentralFreeIntermediate
	case shared.WarehouseCategoryFreeCentraIntermediate:
		return NodeFreeCentralIntermediate
	default:
		return NodeUnrecognized
	}
}
