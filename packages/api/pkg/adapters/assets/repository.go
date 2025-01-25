package assets

import (
	"context"

	"github.com/AssetPortal/assets-api/pkg/model"
)

type Repository interface {
	CreateAsset(ctx context.Context, asset *model.Asset) (*model.Asset, error)
	GetAssetByID(ctx context.Context, id string) (*model.Asset, error)
	GetAssets(ctx context.Context, filters *model.GetAssetsInput) ([]*model.Asset, error)
	UpdateAsset(ctx context.Context, asset *model.Asset) error
	DeleteAsset(ctx context.Context, id, address string) error
}
