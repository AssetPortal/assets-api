package assets

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/uptrace/bun"
)

type AssetsRepository struct {
	db *bun.DB
}

func NewAssetsRepository(db *bun.DB) *AssetsRepository {
	return &AssetsRepository{db: db}
}

func (repo *AssetsRepository) CreateAsset(ctx context.Context, asset *model.Asset) (*model.Asset, error) {
	_, err := repo.db.NewInsert().Model(asset).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to store asset in database: '%s'", err)
	}
	return asset, nil
}
func (repo *AssetsRepository) GetAssetByID(ctx context.Context, id string) (*model.Asset, error) {
	var dbAsset model.Asset
	err := repo.db.NewSelect().Model(&dbAsset).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &dbAsset, fmt.Errorf("failed to query asset by id: '%s'", err)
	}

	return &dbAsset, nil
}
func (repo *AssetsRepository) GetAssets(ctx context.Context, filters *model.GetAssetsInput) ([]*model.Asset, error) {
	var dbAssets []*model.Asset
	query := repo.db.NewSelect().Model(&dbAssets)

	if filters.Address != nil {
		query = query.Where("address = ?", *filters.Address)
	}

	if filters.ID != nil {
		query = query.Where("id = ?", *filters.ID)
	}

	if filters.Order.Order != nil {
		orderClause := fmt.Sprintf("'%s' '%s'", *filters.Order.Order, "ASC")
		if filters.Order.Ascending != nil && !*filters.Order.Ascending {
			orderClause = fmt.Sprintf("'%s' '%s'", *filters.Order.Order, "DESC")
		}
		query = query.Order(orderClause)
	}

	if filters.Pagination.Limit != nil {
		query = query.Limit(*filters.Pagination.Limit)
	}
	if filters.Pagination.Offset != nil {
		query = query.Offset(*filters.Pagination.Offset)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: '%s'", err)
	}

	return dbAssets, nil
}

func (repo *AssetsRepository) UpdateAsset(ctx context.Context, asset *model.Asset) error {
	res, err := repo.db.NewUpdate().Model(asset).Where("id = ? AND address = ?", asset.ID, asset.Address).OmitZero().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update asset in database: '%s'", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update asset in database: '%s'", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("asset with id '%s' and address '%s' does not exist", asset.ID, asset.Address)
	}
	return nil
}
func (repo *AssetsRepository) DeleteAsset(ctx context.Context, id, address string) error {
	res, err := repo.db.NewDelete().Model(&model.Asset{ID: id, Address: address}).Where("id = ? AND address = ?", id, address).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete asset in database: '%s'", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update asset in database: '%s'", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("asset with id '%s' and address '%s' does not exist", id, address)
	}
	return nil
}
