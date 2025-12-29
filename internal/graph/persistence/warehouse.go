package persistence

import (
	"context"

	"github.com/DimKa163/dalty/internal/db"
	"github.com/DimKa163/dalty/internal/graph/core"
)

const (
	GetAllWarehouse = `SELECT nrb_sub_warehouse.id,
       	nrb_sub_warehouse.nrb_fnrec,
       	nrb_name,
       	nrb_is_active,
       	bpm_only_stock_pickup_allowed,
       	nrb_sender_id,
       	nrb_recipient_id,
       	sc.nrb_fnrec,
       	sc.nrb_available_for_balances,
       	nw.id,
       	nw.nrb_fnrec,
       	nw.nrb_address,
       	nw.bpm_descriptor_group_name,
       	tz.id,
       	tz.code
		FROM public.nrb_sub_warehouse
		JOIN nrb_sub_warehouse_categories sc on sc.id=nrb_sub_warehouse.nrb_category_id
		JOIN public.nrb_warehouse nw on nw.id = nrb_sub_warehouse.nrb_warehouse_id
		LEFT JOIN public.time_zone tz on tz.id=nw.ask_time_zone_id
		WHERE nrb_is_active = true`
)

type WarehouseRepository struct {
	db db.QueryExecutor
}

func (w WarehouseRepository) GetAll(ctx context.Context) ([]*core.Warehouse, error) {
	var warehouses []*core.Warehouse
	rows, err := w.db.Query(ctx, GetAllWarehouse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var warehouse core.Warehouse
		if err := warehouse.Scan(rows); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}
	return warehouses, nil
}

func NewWarehouseRepository(db db.QueryExecutor) *WarehouseRepository {
	return &WarehouseRepository{
		db: db,
	}
}
