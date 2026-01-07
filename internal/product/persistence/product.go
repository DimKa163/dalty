package persistence

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DimKa163/dalty/pkg/daltyerrors"

	"github.com/DimKa163/dalty/internal/db"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
)

const (
	GetByIDStmt = `SELECT 
	id,
    name,
    type_id,
    nrb_type_production_id,
    smr_fnrec,
    is_archive,
    nrb_integration_id,
    smr_is_service,
    smr_product_group_flag_id,
    category_id,
    smr_series_id,
    nrb_account_product_id,
    ask_non_standart_category_id,
    nrb_count_mv,
    ask_pack_volume,
    ask_pack_length,
    ask_pack_width,
    ask_pack_height,
    ask_weight 
	FROM public.product
    WHERE id=$1`
	GetByFnrecStmt = `SELECT 
	id,
    name,
    type_id,
    nrb_type_production_id,
    smr_fnrec,
    is_archive,
    nrb_integration_id,
    smr_is_service,
    smr_product_group_flag_id,
    category_id,
    smr_series_id,
    nrb_account_product_id,
    ask_non_standart_category_id,
    nrb_count_mv,
    ask_pack_volume,
    ask_pack_length,
    ask_pack_width,
    ask_pack_height,
    ask_weight 
	FROM public.product
    WHERE smr_fnrec=$1`
	GetByIntegrationIDStmt = `SELECT 
	id,
    name,
    type_id,
    nrb_type_production_id,
    smr_fnrec,
    is_archive,
    nrb_integration_id,
    smr_is_service,
    smr_product_group_flag_id,
    category_id,
    smr_series_id,
    nrb_account_product_id,
    ask_non_standart_category_id,
    nrb_count_mv,
    ask_pack_volume,
    ask_pack_length,
    ask_pack_width,
    ask_pack_height,
    ask_weight 
	FROM public.product
    WHERE nrb_integration_id=$1`
)

type ProductRepository struct {
	db db.QueryExecutor
}

func NewProductRepository(db db.QueryExecutor) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(ctx context.Context, id guid.Guid) (*core.Product, error) {
	row := r.db.QueryRow(ctx, GetByIDStmt, id)
	prd, err := mapProduct(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, daltyerrors.NewNotFoundError(daltyerrors.ErrNotFound, "product not found by id", id)
		}
		return nil, err
	}
	return prd, nil
}

func (r *ProductRepository) GetByFnrec(ctx context.Context, fnrec string) (*core.Product, error) {
	row := r.db.QueryRow(ctx, GetByFnrecStmt, fnrec)
	prd, err := mapProduct(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, daltyerrors.NewNotFoundError(daltyerrors.ErrNotFound, "product not found by fnrec", fnrec)
		}
		return nil, err
	}
	return prd, nil
}

func (r *ProductRepository) GetByIntegrationID(ctx context.Context, integrationID string) (*core.Product, error) {
	row := r.db.QueryRow(ctx, GetByIntegrationIDStmt, integrationID)
	prd, err := mapProduct(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, daltyerrors.NewNotFoundError(daltyerrors.ErrNotFound, "product not found by integrationID", integrationID)
		}
		return nil, err
	}
	return prd, nil
}

func mapProduct(row pgx.Row) (*core.Product, error) {
	var product core.Product
	var productID guid.Guid
	var name string
	var typeID string
	var typeProductionID string
	var fnrec string
	var isArchive bool
	var integrationID string
	var isService bool
	var productGroupFlagID sql.NullString
	var categoryID sql.NullString
	var seriesID sql.NullString
	var accountProviderID sql.NullString
	var standardCategory sql.NullString
	var countMa int
	var pack core.ProductPack

	if err := row.Scan(&productID,
		&name,
		&typeID,
		&typeProductionID,
		&fnrec,
		&isArchive,
		&integrationID,
		&isService,
		&productGroupFlagID,
		&categoryID,
		&seriesID,
		&accountProviderID,
		&standardCategory,
		&countMa,
		&pack.Volume,
		&pack.Length,
		&pack.Width,
		&pack.Height,
		&pack.Weight); err != nil {
		return nil, err
	}
	product.ID = productID
	product.Name = name
	product.Fnrec = fnrec
	product.IsArchive = isArchive
	product.IsService = isService
	product.IntegrationID = integrationID
	product.CountMa = countMa
	product.ProductPack = pack
	product.Type = core.ProductType(typeID)
	product.ProductionType = core.ProductionType(typeProductionID)
	if productGroupFlagID.Valid {
		product.Group = core.ProductGroup(productGroupFlagID.String)
	}
	if categoryID.Valid {
		product.CategoryID = categoryID.String
	}
	if seriesID.Valid {
		product.SeriesID = seriesID.String
	}
	if accountProviderID.Valid {
		product.AccountProviderId = accountProviderID.String
	}
	if standardCategory.Valid {
		product.NonStandardCategory = standardCategory.String
	}
	return &product, nil
}
