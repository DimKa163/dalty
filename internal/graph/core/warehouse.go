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
	Category               *WarehouseCategory
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
	w.Category = &WarehouseCategory{
		Fnrec:               categoryFnrec.String,
		AvailableForBalance: availableForBalances,
	}
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

type WarehouseCategory struct {
	Fnrec               string
	AvailableForBalance bool
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
