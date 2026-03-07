package core

import (
	"context"
	"database/sql"

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
	Type                   WarehouseCategory
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
	w.Type = WarehouseCategory(categoryFnrec.String)
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
type WarehouseCategory string

const (
	WarehouseCategoryMall                    WarehouseCategory = "80010000004322B0"
	WarehouseCategoryCentral                 WarehouseCategory = "80010000004322AC"
	WarehouseCategoryFree                    WarehouseCategory = "800100000044849C"
	WarehouseCategoryMain                    WarehouseCategory = "80010000004322AE"
	WarehouseCategoryTransit                 WarehouseCategory = "800100000044849B"
	WarehouseCategoryReservation             WarehouseCategory = "800100000044A9C9"
	WarehouseCategoryLoses                   WarehouseCategory = "8001000000449309"
	WarehouseCategoryMarketing               WarehouseCategory = "80010000004490A3"
	WarehouseCategoryExposition              WarehouseCategory = "8001000000449051"
	WarehouseCategoryPartner                 WarehouseCategory = "800100000044FD7E"
	WarehouseCategoryPartner2                WarehouseCategory = "800100000045D682"
	WarehouseCategoryFree2                   WarehouseCategory = "800100000045E480"
	WarehouseCategoryProblem                 WarehouseCategory = "800100000045EF15"
	WarehouseCategoryRefund                  WarehouseCategory = "80010000004613D1"
	WarehouseCategoryProduction              WarehouseCategory = "80010000004613D2"
	WarehouseCategoryRecycling               WarehouseCategory = "8001000000461422"
	WarehouseCategoryService                 WarehouseCategory = "8001000000461424"
	WarehouseCategoryMaterial                WarehouseCategory = "8001000000461588"
	WarehouseCategoryMarkdown                WarehouseCategory = "80010000004615E6"
	WarehouseCategoryBuffer                  WarehouseCategory = "80010000004616F6"
	WarehouseCategoryDiscount                WarehouseCategory = "8001000000432AE5"
	WarehouseCategoryCentralMainIntermediate WarehouseCategory = "80010000004322AD"
	WarehouseCategoryMainCentraIntermediate  WarehouseCategory = "80010000004322AF"
	WarehouseCategoryCentraFreeIntermediate  WarehouseCategory = "8001000000448BF4"
	WarehouseCategoryFreeCentraIntermediate  WarehouseCategory = "8001000000464D17"
)
