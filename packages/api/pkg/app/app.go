package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/AssetPortal/assets-api/pkg/adapters/assets"
	"github.com/AssetPortal/assets-api/pkg/adapters/storage"
	"github.com/AssetPortal/assets-api/pkg/adapters/tokens"
	"github.com/AssetPortal/assets-api/pkg/config"
	appError "github.com/AssetPortal/assets-api/pkg/error"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type AssetsApp struct {
	cfg              *config.Configuration
	db               *bun.DB
	tokensRepository tokens.Repository
	assetsRepository assets.Repository
	storageClient    storage.Client
	log              *logrus.Logger
}

func NewAssetsApp(
	cfg *config.Configuration,
	db *bun.DB,
	tokensRepository tokens.Repository,
	assetsRepository assets.Repository,
	storageClient storage.Client,
	log *logrus.Logger,
) *AssetsApp {
	return &AssetsApp{
		cfg:              cfg,
		db:               db,
		tokensRepository: tokensRepository,
		assetsRepository: assetsRepository,
		storageClient:    storageClient,
		log:              log,
	}
}

func (app *AssetsApp) CreateToken(ctx context.Context) (*model.Token, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		app.log.Errorf("error generating random token: '%s'", err)
		return nil, appError.ErrGeneratingToken
	}
	now := time.Now()
	newToken := &model.Token{
		Token:     hex.EncodeToString(tokenBytes),
		CreatedAt: now,
		ExpiresAt: now.Add(app.cfg.TokenExpiration),
		Used:      false,
	}
	token, err := app.tokensRepository.CreateToken(ctx, newToken)
	if err != nil {
		app.log.Errorf("error creating token: '%s'", err)
		return nil, appError.ErrCreatingToken
	}
	return token, nil
}

func (app *AssetsApp) CreateAsset(ctx context.Context, asset *model.Asset) (*model.Asset, error) {
	token, err := app.assetsRepository.CreateAsset(ctx, asset)
	if err != nil {
		if strings.Contains(err.Error(), "assets_id_key") {
			return nil, appError.ErrCreatingAssetIDExists
		}
		app.log.Errorf("error creating asset: '%s'", err)
		return nil, appError.ErrCreatingAsset
	}
	return token, nil
}

func (app *AssetsApp) GetAssetByID(ctx context.Context, id string) (*model.Asset, error) {
	asset, err := app.assetsRepository.GetAssetByID(ctx, id)
	if err != nil {
		app.log.Errorf("error getting asset by id '%s': '%s'", id, err)
		return nil, appError.ErrGettingAsset
	}
	return asset, nil
}

func (app *AssetsApp) GetAssets(ctx context.Context, filters *model.GetAssetsInput) ([]*model.Asset, error) {
	asset, err := app.assetsRepository.GetAssets(ctx, filters)
	if err != nil {
		app.log.Errorf("error getting assets: '%s'", err)
		return nil, appError.ErrGettingAsset
	}
	return asset, nil
}

func (app *AssetsApp) UploadFile(ctx context.Context, fileKey string, fileBytes []byte, contentType string) (*model.URL, error) {
	url, err := app.storageClient.UploadFile(ctx, fileKey, fileBytes, contentType)
	if err != nil {
		app.log.Errorf("error uploading file: '%s'", err)
		return nil, appError.ErrUploadingFile
	}
	return &model.URL{
		URL: url,
	}, nil
}

func (app *AssetsApp) UpdateAsset(ctx context.Context, asset *model.Asset) error {
	err := app.assetsRepository.UpdateAsset(ctx, asset)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return appError.ErrAssetDoesNotBelongToTheUser
		}
		app.log.Errorf("error updating asset by id '%s': '%s'", asset.ID, err)
		return appError.ErrUpdatingAsset
	}
	return nil
}

func (app *AssetsApp) DeleteAsset(ctx context.Context, id, address string) error {
	err := app.assetsRepository.DeleteAsset(ctx, id, address)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return appError.ErrAssetDoesNotBelongToTheUser
		}
		app.log.Errorf("error deleting asset by id '%s' and address '%s': '%s'", id, address, err)
		return appError.ErrDeletingAsset
	}
	return nil
}

func (app *AssetsApp) Config() *config.Configuration {
	return app.cfg
}
