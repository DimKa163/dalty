package persistence

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DimKa163/dalty/internal/db"
	"github.com/DimKa163/dalty/internal/rest/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

const (
	findRestStmt = `SELECT id,
		nrb_product_id,
		nrb_sub_warehouse_id,
		bpm_filial_id,
		nrb_integration_id,
		nrb_quantity
	FROM public.nrb_product_balances
	WHERE nrb_product_id = $1
	AND nrb_sub_warehouse_id = $2
	AND bpm_filial_id = $3`
)

type RestRepository struct {
	db db.QueryExecutor
}

func NewRestRepository(db db.QueryExecutor) *RestRepository {
	return &RestRepository{db: db}
}

func (r *RestRepository) FindRest(ctx context.Context, productID, subWarehouse, filialID guid.Guid) (*core.Rest, error) {
	var rest core.Rest
	var id guid.Guid
	var prdId guid.Guid
	var subWhId guid.Guid
	var filWhId guid.Guid
	var integrationId string
	var quantity sql.NullString
	if err := r.db.QueryRow(ctx, findRestStmt, productID, subWarehouse, filialID).Scan(&id, &prdId, &subWhId, &filWhId, &integrationId, &quantity); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, daltyerrors.NewNotFoundError(daltyerrors.ErrNotFound, "rests not found", id)
		}
		return nil, err
	}
	rest.ID = id
	rest.ProductID = prdId
	rest.SubWarehouseID = subWhId
	rest.FilialID = filWhId
	rest.IntegrationID = integrationId
	if quantity.Valid {
		d, err := decimal.NewFromString(quantity.String)
		if err != nil {
			return nil, err
		}
		rest.Quantity = d
	}
	return &rest, nil
}
