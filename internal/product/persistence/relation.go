package persistence

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DimKa163/dalty/internal/db"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltymodel"
	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
)

const (
	GetByLeftIDStmt = `SELECT
		nrb_related_product.id,
		nrb_product_sku_id,
		nrb_amount_mv,
		p1.id,
    	p1.name,
    	p1.type_id,
    	p1.nrb_type_production_id,
    	p1.smr_fnrec,
    	p1.is_archive,
    	p1.nrb_integration_id,
    	p1.smr_is_service,
    	p1.smr_product_group_flag_id,
    	p1.category_id,
    	p1.smr_series_id,
    	p1.nrb_account_product_id,
    	p1.ask_non_standart_category_id,
    	p1.nrb_count_mv,
    	p1.ask_pack_volume,
    	p1.ask_pack_length,
    	p1.ask_pack_width,
    	p1.ask_pack_height,
    	p1.ask_weight
	FROM public.nrb_related_product
	JOIN public.product p on nrb_related_product.nrb_product_sku_id = p.id
	JOIN public.product p1 on nrb_related_product.nrb_product_mv_id = p1.id
	WHERE p.id=$1`
	GetByLeftFnrecStmt = `SELECT
		nrb_related_product.id,
		nrb_product_sku_id,
		p1.id,
    	p1.name,
    	p1.type_id,
    	p1.nrb_type_production_id,
    	p1.smr_fnrec,
    	p1.is_archive,
    	p1.nrb_integration_id,
    	p1.smr_is_service,
    	p1.smr_product_group_flag_id,
    	p1.category_id,
    	p1.smr_series_id,
    	p1.nrb_account_product_id,
    	p1.ask_non_standart_category_id,
    	p1.nrb_count_mv,
    	p1.ask_pack_volume,
    	p1.ask_pack_length,
    	p1.ask_pack_width,
    	p1.ask_pack_height,
    	p1.ask_weight,
		nrb_amount_mv
	FROM public.nrb_related_product
	JOIN public.product p on nrb_related_product.nrb_product_sku_id = p.id
	JOIN public.product p1 on nrb_related_product.nrb_product_mv_id = p1.id
	WHERE p.smr_fnrec=$1`
	GetByLeftIntegrationIDStmt = `SELECT
		nrb_related_product.id,
		nrb_product_sku_id,
		p1.id,
    	p1.name,
    	p1.type_id,
    	p1.nrb_type_production_id,
    	p1.smr_fnrec,
    	p1.is_archive,
    	p1.nrb_integration_id,
    	p1.smr_is_service,
    	p1.smr_product_group_flag_id,
    	p1.category_id,
    	p1.smr_series_id,
    	p1.nrb_account_product_id,
    	p1.ask_non_standart_category_id,
    	p1.nrb_count_mv,
    	p1.ask_pack_volume,
    	p1.ask_pack_length,
    	p1.ask_pack_width,
    	p1.ask_pack_height,
    	p1.ask_weight,
		nrb_amount_mv
	FROM public.nrb_related_product
	JOIN public.product p on nrb_related_product.nrb_product_sku_id = p.id
	JOIN public.product p1 on nrb_related_product.nrb_product_mv_id = p1.id
	WHERE p.nrb_integration_id=$1`
	GetByRightIDStmt = `SELECT 
    	l.id, 
    	l.nrb_product_sku_id, 
    	l.nrb_product_mv_id, 
    	l.nrb_amount_mv, 
    	r.id, 
    	r.nrb_product_sku_id, 
    	r.nrb_product_mv_id, 
    	r.nrb_amount_mv,
    	p.id,
    	p.name,
    	p.type_id,
    	p.nrb_type_production_id,
    	p.smr_fnrec,
    	p.is_archive,
    	p.nrb_integration_id,
    	p.smr_is_service,
    	p.smr_product_group_flag_id,
    	p.category_id,
    	p.smr_series_id,
    	p.nrb_account_product_id,
    	p.ask_non_standart_category_id,
    	p.nrb_count_mv,
    	p.ask_pack_volume,
    	p.ask_pack_length,
    	p.ask_pack_width,
    	p.ask_pack_height,
    	p.ask_weight
	FROM nrb_related_product AS l
 	JOIN nrb_related_product AS r ON r.nrb_product_sku_id = l.nrb_product_sku_id
	JOIN product p ON l.nrb_product_sku_id = p.id
	WHERE l.nrb_product_mv_id = $1
  	AND r.nrb_product_mv_id=$2`
)

type RelationRepository struct {
	db db.QueryExecutor
}

func NewRelationRepository(db db.QueryExecutor) *RelationRepository {
	return &RelationRepository{db: db}
}

func (r *RelationRepository) GetByLeftID(ctx context.Context, id guid.Guid) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftIDStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readDirectRelation(rows)
}
func (r *RelationRepository) GetByLeftFnrec(ctx context.Context, fnrec string) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftFnrecStmt, fnrec)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readDirectRelation(rows)
}

func (r *RelationRepository) GetByLeftIntegrationID(ctx context.Context, integrationID string) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftIntegrationIDStmt, integrationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readDirectRelation(rows)
}

func (r *RelationRepository) GetByRightID(ctx context.Context, lid, rid guid.Guid) (*core.Relation, *core.Relation, error) {
	var lrelationID guid.Guid
	var lParentProductID guid.Guid
	var lSubProductID guid.Guid
	var lAmount int32
	var rrelationID guid.Guid
	var rParentProductID guid.Guid
	var rSubProductID guid.Guid
	var rAmount int32
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
	var countMa int32
	var volume float64
	var length float64
	var width float64
	var height float64
	var weight float64
	if err := r.db.QueryRow(ctx, GetByRightIDStmt, lid, rid).Scan(
		&lrelationID,
		&lParentProductID,
		&lSubProductID,
		&lAmount,
		&rrelationID,
		&rParentProductID,
		&rSubProductID,
		&rAmount,
		&productID,
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
		&volume,
		&length,
		&width,
		&height,
		&weight); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, daltyerrors.NewNotFoundError(daltyerrors.ErrNotFound, "product net not found", lid, rid)
		}
		return nil, nil, err

	}
	var lr core.Relation
	var rr core.Relation
	lr.ID = lrelationID
	lr.LeftID = lParentProductID
	lr.RightID = lSubProductID
	lr.Amount = lAmount
	lr.Left = &product
	rr.ID = rrelationID
	rr.LeftID = rParentProductID
	rr.RightID = rSubProductID
	rr.Amount = rAmount
	rr.Left = &product
	product.ID = productID
	product.Name = name
	product.Fnrec = fnrec
	product.IsArchive = isArchive
	product.IsService = isService
	product.IntegrationID = integrationID
	product.CountMa = countMa
	product.Volume = volume
	product.Length = length
	product.Width = width
	product.Height = height
	product.Weight = weight
	product.Type = daltymodel.ProductType(typeID)
	product.ProductionType = daltymodel.ProductionType(typeProductionID)
	if productGroupFlagID.Valid {
		product.Group = daltymodel.ProductGroup(productGroupFlagID.String)
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
	return &lr, &rr, nil
}

func readDirectRelation(rows pgx.Rows) ([]*core.Relation, error) {
	relations := make([]*core.Relation, 0)
	for rows.Next() {
		var id guid.Guid
		var left guid.Guid
		var amount int32
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
		var countMa int32
		var volume float64
		var length float64
		var width float64
		var height float64
		var weight float64
		if err := rows.Scan(&id,
			&left,
			&amount,
			&productID,
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
			&volume,
			&length,
			&width,
			&height,
			&weight); err != nil {
			return nil, err
		}
		var relation core.Relation
		relation.ID = id
		relation.Amount = amount
		relation.LeftID = left
		relation.RightID = productID
		relation.Right = &core.Product{}
		relation.Right.ID = productID
		relation.Right.Name = name
		relation.Right.Fnrec = fnrec
		relation.Right.IsArchive = isArchive
		relation.Right.IsService = isService
		relation.Right.IntegrationID = integrationID
		relation.Right.CountMa = countMa
		relation.Right.Volume = volume
		relation.Right.Length = length
		relation.Right.Width = width
		relation.Right.Height = height
		relation.Right.Weight = weight
		relation.Right.Type = daltymodel.ProductType(typeID)
		relation.Right.ProductionType = daltymodel.ProductionType(typeProductionID)
		if productGroupFlagID.Valid {
			relation.Right.Group = daltymodel.ProductGroup(productGroupFlagID.String)
		}
		if categoryID.Valid {
			relation.Right.CategoryID = categoryID.String
		}
		if seriesID.Valid {
			relation.Right.SeriesID = seriesID.String
		}
		if accountProviderID.Valid {
			relation.Right.AccountProviderId = accountProviderID.String
		}
		if standardCategory.Valid {
			relation.Right.NonStandardCategory = standardCategory.String
		}
		relations = append(relations, &relation)

	}
	return relations, nil
}
