package persistence

import (
	"context"
	"database/sql"

	"github.com/DimKa163/dalty/internal/db"
	"github.com/DimKa163/dalty/internal/product/core"
	"github.com/beevik/guid"
	"github.com/jackc/pgx/v5"
)

const (
	GetCountStmt       = `SELECT COUNT(*) FROM public.nrb_related_product`
	GetAllRelationStmt = `SELECT 
    	id,
    	nrb_amount_mv, 
    	nrb_product_sku_id, 
    	nrb_product_mv_id
		FROM public.nrb_related_product`
	GetByLeftByIDStmt = `SELECT
		 rel.id,
    rel.nrb_amount_mv,
    l.id,
    l.name,
    l.type_id,
    l.nrb_type_production_id,
    l.smr_fnrec,
    l.is_archive,
    l.nrb_integration_id,
    l.smr_is_service,
    l.smr_product_group_flag_id,
    l.category_id,
    l.smr_series_id,
    l.nrb_account_product_id,
    l.ask_non_standart_category_id,
    l.nrb_count_mv,
    l.ask_pack_volume,
    l.ask_pack_length,
    l.ask_pack_width,
    l.ask_pack_height,
    l.ask_weight,
    r.id,
    r.name,
    r.type_id,
    r.nrb_type_production_id,
    r.smr_fnrec,
    r.is_archive,
    r.nrb_integration_id,
    r.smr_is_service,
    r.smr_product_group_flag_id,
    r.category_id,
    r.smr_series_id,
    r.nrb_account_product_id,
    r.ask_non_standart_category_id,
    r.nrb_count_mv,
    r.ask_pack_volume,
    r.ask_pack_length,
    r.ask_pack_width,
    r.ask_pack_height,
    r.ask_weight
	FROM public.nrb_related_product rel
	JOIN public.product l ON l.id = rel.nrb_product_sku_id
	JOIN public.product r ON r.id = rel.nrb_product_mv_id
	WHERE l.id = $1`
	GetByLeftFnrecStmt = `SELECT
		 rel.id,
    rel.nrb_amount_mv,
    l.id,
    l.name,
    l.type_id,
    l.nrb_type_production_id,
    l.smr_fnrec,
    l.is_archive,
    l.nrb_integration_id,
    l.smr_is_service,
    l.smr_product_group_flag_id,
    l.category_id,
    l.smr_series_id,
    l.nrb_account_product_id,
    l.ask_non_standart_category_id,
    l.nrb_count_mv,
    l.ask_pack_volume,
    l.ask_pack_length,
    l.ask_pack_width,
    l.ask_pack_height,
    l.ask_weight,
    r.id,
    r.name,
    r.type_id,
    r.nrb_type_production_id,
    r.smr_fnrec,
    r.is_archive,
    r.nrb_integration_id,
    r.smr_is_service,
    r.smr_product_group_flag_id,
    r.category_id,
    r.smr_series_id,
    r.nrb_account_product_id,
    r.ask_non_standart_category_id,
    r.nrb_count_mv,
    r.ask_pack_volume,
    r.ask_pack_length,
    r.ask_pack_width,
    r.ask_pack_height,
    r.ask_weight
	FROM public.nrb_related_product rel
	JOIN public.product l ON l.id = rel.nrb_product_sku_id
	JOIN public.product r ON r.id = rel.nrb_product_mv_id
	WHERE l.smr_fnrec = $1`
	GetByLeftIntegrationIDStmt = `SELECT
    	    rel.id,
    rel.nrb_amount_mv,
    l.id,
    l.name,
    l.type_id,
    l.nrb_type_production_id,
    l.smr_fnrec,
    l.is_archive,
    l.nrb_integration_id,
    l.smr_is_service,
    l.smr_product_group_flag_id,
    l.category_id,
    l.smr_series_id,
    l.nrb_account_product_id,
    l.ask_non_standart_category_id,
    l.nrb_count_mv,
    l.ask_pack_volume,
    l.ask_pack_length,
    l.ask_pack_width,
    l.ask_pack_height,
    l.ask_weight,
    r.id,
    r.name,
    r.type_id,
    r.nrb_type_production_id,
    r.smr_fnrec,
    r.is_archive,
    r.nrb_integration_id,
    r.smr_is_service,
    r.smr_product_group_flag_id,
    r.category_id,
    r.smr_series_id,
    r.nrb_account_product_id,
    r.ask_non_standart_category_id,
    r.nrb_count_mv,
    r.ask_pack_volume,
    r.ask_pack_length,
    r.ask_pack_width,
    r.ask_pack_height,
    r.ask_weight
	FROM public.nrb_related_product rel
	JOIN public.product l ON l.id = rel.nrb_product_sku_id
	JOIN public.product r ON r.id = rel.nrb_product_mv_id
	WHERE l.nrb_integration_id = $1`
	GetByRightIDStmt = `SELECT
    	    rel.id,
    rel.nrb_amount_mv,
    l.id,
    l.name,
    l.type_id,
    l.nrb_type_production_id,
    l.smr_fnrec,
    l.is_archive,
    l.nrb_integration_id,
    l.smr_is_service,
    l.smr_product_group_flag_id,
    l.category_id,
    l.smr_series_id,
    l.nrb_account_product_id,
    l.ask_non_standart_category_id,
    l.nrb_count_mv,
    l.ask_pack_volume,
    l.ask_pack_length,
    l.ask_pack_width,
    l.ask_pack_height,
    l.ask_weight,
    r.id,
    r.name,
    r.type_id,
    r.nrb_type_production_id,
    r.smr_fnrec,
    r.is_archive,
    r.nrb_integration_id,
    r.smr_is_service,
    r.smr_product_group_flag_id,
    r.category_id,
    r.smr_series_id,
    r.nrb_account_product_id,
    r.ask_non_standart_category_id,
    r.nrb_count_mv,
    r.ask_pack_volume,
    r.ask_pack_length,
    r.ask_pack_width,
    r.ask_pack_height,
    r.ask_weight
	FROM public.nrb_related_product rel
	JOIN public.product l ON l.id = rel.nrb_product_sku_id
	JOIN public.product r ON r.id = rel.nrb_product_mv_id
	WHERE r.id = $1 AND rel.nrb_amount_mv=$2`
)

type RelationRepository struct {
	db db.QueryExecutor
}

func NewRelationRepository(db db.QueryExecutor) *RelationRepository {
	return &RelationRepository{db: db}
}

func (r *RelationRepository) GetAll(ctx context.Context) ([]*core.Relation, error) {
	var count int
	if err := r.db.QueryRow(ctx, GetCountStmt).Scan(&count); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, GetAllRelationStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var index int
	relations := make([]*core.Relation, count)
	for rows.Next() {
		var id guid.Guid
		var amount int
		var left guid.Guid
		var right guid.Guid
		if err := rows.Scan(&id, &amount, &left, &right); err != nil {
			return nil, err
		}
		relations[index] = &core.Relation{
			ID:      id,
			Amount:  amount,
			LeftID:  left,
			RightID: right,
		}
		index++
	}
	return relations, nil
}

func (r *RelationRepository) GetByLeftID(ctx context.Context, id guid.Guid) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftByIDStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRelation(rows)
}
func (r *RelationRepository) GetByLeftFnrec(ctx context.Context, fnrec string) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftFnrecStmt, fnrec)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRelation(rows)
}

func (r *RelationRepository) GetByLeftIntegrationID(ctx context.Context, integrationID string) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByLeftIntegrationIDStmt, integrationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRelation(rows)
}

func (r *RelationRepository) GetByRightID(ctx context.Context, id guid.Guid, amount int) ([]*core.Relation, error) {
	rows, err := r.db.Query(ctx, GetByRightIDStmt, id, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readRelation(rows)
}

func readRelation(rows pgx.Rows) ([]*core.Relation, error) {
	relations := make([]*core.Relation, 0)
	for rows.Next() {
		var id guid.Guid
		var amount int
		var left guid.Guid
		var right guid.Guid
		var leftName string
		var leftTypeID string
		var leftTypeProductionID string
		var leftFnrec string
		var leftIsArchive bool
		var leftIntegrationID string
		var leftIsService bool
		var leftProductGroupFlagID sql.NullString
		var leftCategoryID sql.NullString
		var leftSeriesID sql.NullString
		var leftAccountProviderID sql.NullString
		var leftStandardCategory sql.NullString
		var leftCountMa int
		var leftPack core.ProductPack
		var rightName string
		var rightTypeID string
		var rightTypeProductionID string
		var rightFnrec string
		var rightIsArchive bool
		var rightIntegrationID string
		var rightIsService bool
		var rightProductGroupFlagID sql.NullString
		var rightCategoryID sql.NullString
		var rightSeriesID sql.NullString
		var rightAccountProviderID sql.NullString
		var rightStandardCategory sql.NullString
		var rightCountMa int
		var rightPack core.ProductPack
		if err := rows.Scan(&id,
			&amount,
			&left,
			&leftName,
			&leftTypeID,
			&leftTypeProductionID,
			&leftFnrec,
			&leftIsArchive,
			&leftIntegrationID,
			&leftIsService,
			&leftProductGroupFlagID,
			&leftCategoryID,
			&leftSeriesID,
			&leftAccountProviderID,
			&leftStandardCategory,
			&leftCountMa,
			&leftPack.Volume,
			&leftPack.Length,
			&leftPack.Width,
			&leftPack.Height,
			&leftPack.Weight,
			&right,
			&rightName,
			&rightTypeID,
			&rightTypeProductionID,
			&rightFnrec,
			&rightIsArchive,
			&rightIntegrationID,
			&rightIsService,
			&rightProductGroupFlagID,
			&rightCategoryID,
			&rightSeriesID,
			&rightAccountProviderID,
			&rightStandardCategory,
			&rightCountMa,
			&rightPack.Volume,
			&rightPack.Length,
			&rightPack.Width,
			&rightPack.Height,
			&rightPack.Weight); err != nil {
			return nil, err
		}
		var relation core.Relation
		relation.ID = id
		relation.Amount = amount
		relation.LeftID = left
		relation.RightID = right
		relation.Left = core.Product{
			ID:            left,
			Name:          leftName,
			Fnrec:         leftFnrec,
			IsArchive:     leftIsArchive,
			IsService:     leftIsService,
			IntegrationID: leftIntegrationID,
			CountMa:       leftCountMa,
			ProductPack:   leftPack,
		}
		relation.Left.Type = core.ProductType(leftTypeID)
		relation.Left.ProductionType = core.ProductionType(leftTypeProductionID)
		if leftProductGroupFlagID.Valid {
			relation.Left.Group = core.ProductGroup(leftProductGroupFlagID.String)
		}
		if leftCategoryID.Valid {
			relation.Left.CategoryID = leftCategoryID.String
		}
		if leftSeriesID.Valid {
			relation.Left.SeriesID = leftSeriesID.String
		}
		if leftAccountProviderID.Valid {
			relation.Left.AccountProviderId = leftAccountProviderID.String
		}
		if leftStandardCategory.Valid {
			relation.Left.NonStandardCategory = leftStandardCategory.String
		}
		relation.Right = core.Product{
			ID:            right,
			Name:          rightName,
			Fnrec:         rightFnrec,
			IsArchive:     rightIsArchive,
			IsService:     rightIsService,
			IntegrationID: rightIntegrationID,
			CountMa:       rightCountMa,
			ProductPack:   rightPack,
		}
		relation.Right.Type = core.ProductType(rightTypeID)
		relation.Right.ProductionType = core.ProductionType(rightTypeProductionID)
		if rightProductGroupFlagID.Valid {
			relation.Right.Group = core.ProductGroup(rightProductGroupFlagID.String)
		}
		if rightCategoryID.Valid {
			relation.Right.CategoryID = rightCategoryID.String
		}
		if rightSeriesID.Valid {
			relation.Right.SeriesID = rightSeriesID.String
		}
		if rightAccountProviderID.Valid {
			relation.Right.AccountProviderId = rightAccountProviderID.String
		}
		if rightStandardCategory.Valid {
			relation.Right.NonStandardCategory = rightStandardCategory.String
		}
		relations = append(relations, &relation)

	}
	return relations, nil
}
