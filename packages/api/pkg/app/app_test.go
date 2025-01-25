package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	assetsMock "github.com/AssetPortal/assets-api/pkg/adapters/assets/mocks"
	tokensMock "github.com/AssetPortal/assets-api/pkg/adapters/tokens/mocks"
	"github.com/AssetPortal/assets-api/pkg/app"
	"github.com/AssetPortal/assets-api/pkg/config"
	appError "github.com/AssetPortal/assets-api/pkg/error"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssetsApp_CreateToken_Success(t *testing.T) {
	mockTokensRepository := new(tokensMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{TokenExpiration: 1 * time.Hour},
		nil,
		mockTokensRepository,
		nil,
		nil,
		mockLogger,
	)

	expectedToken := &model.Token{
		Token:     "mockedToken123",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
		Used:      false,
	}

	mockTokensRepository.On("CreateToken", mock.Anything, mock.AnythingOfType("*model.Token")).
		Return(expectedToken, nil).Once()

	token, err := app.CreateToken(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken.Token, token.Token)

	mockTokensRepository.AssertExpectations(t)
}

func TestAssetsApp_CreateToken_Failure(t *testing.T) {
	mockTokensRepository := new(tokensMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{TokenExpiration: 1 * time.Hour},
		nil,
		mockTokensRepository,
		nil,
		nil,
		mockLogger,
	)

	mockTokensRepository.On("CreateToken", mock.Anything, mock.AnythingOfType("*model.Token")).
		Return(nil, fmt.Errorf("database error")).Once()

	token, err := app.CreateToken(context.Background())

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, appError.ErrCreatingToken, err)

	mockTokensRepository.AssertExpectations(t)
}

func TestAssetsApp_CreateToken_Expiration(t *testing.T) {
	mockTokensRepository := new(tokensMock.Repository)
	mockLogger := logrus.New()

	expirationDuration := 2 * time.Hour
	app := app.NewAssetsApp(
		&config.Configuration{TokenExpiration: expirationDuration},
		nil,
		mockTokensRepository,
		nil,
		nil,
		mockLogger,
	)

	expectedToken := &model.Token{
		Token:     "mockedToken123",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expirationDuration),
		Used:      false,
	}

	mockTokensRepository.On("CreateToken", mock.Anything, mock.AnythingOfType("*model.Token")).
		Return(expectedToken, nil).Once()

	token, err := app.CreateToken(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken.ExpiresAt, token.ExpiresAt)

	mockTokensRepository.AssertExpectations(t)
}
func TestAssetsApp_CreateToken_LongExpiration(t *testing.T) {
	mockTokensRepository := new(tokensMock.Repository)
	mockLogger := logrus.New()

	expirationDuration := 10 * 365 * 24 * time.Hour
	app := app.NewAssetsApp(
		&config.Configuration{TokenExpiration: expirationDuration},
		nil,
		mockTokensRepository,
		nil,
		nil,
		mockLogger,
	)

	expectedToken := &model.Token{
		Token:     "mockedToken123",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expirationDuration),
		Used:      false,
	}

	mockTokensRepository.On("CreateToken", mock.Anything, mock.AnythingOfType("*model.Token")).
		Return(expectedToken, nil).Once()

	token, err := app.CreateToken(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken.ExpiresAt, token.ExpiresAt)

	mockTokensRepository.AssertExpectations(t)
}
func TestAssetsApp_CreateAsset_Success(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	expectedAsset := &model.Asset{
		ID: "mockedAsset123",
	}

	mockAssetsRepository.On("CreateAsset", mock.Anything, mock.AnythingOfType("*model.Asset")).
		Return(expectedAsset, nil).Once()

	asset, err := app.CreateAsset(context.Background(), expectedAsset)

	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset.ID, asset.ID)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_CreateAsset_AssetIDExists(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	expectedAsset := &model.Asset{
		ID: "mockedAsset123",
	}

	mockAssetsRepository.On("CreateAsset", mock.Anything, mock.AnythingOfType("*model.Asset")).
		Return(nil, fmt.Errorf("duplicate key value violates unique constraint \"assets_id_key\"")).Once()

	asset, err := app.CreateAsset(context.Background(), expectedAsset)

	assert.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, appError.ErrCreatingAssetIDExists, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_CreateAsset_Failure(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	expectedAsset := &model.Asset{
		ID: "mockedAsset123",
	}

	mockAssetsRepository.On("CreateAsset", mock.Anything, mock.AnythingOfType("*model.Asset")).
		Return(nil, fmt.Errorf("some unknown error")).Once()

	asset, err := app.CreateAsset(context.Background(), expectedAsset)

	assert.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, appError.ErrCreatingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_CreateAsset_InvalidAsset(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	invalidAsset := &model.Asset{
		ID: "mockedAsset123",
	}

	mockAssetsRepository.On("CreateAsset", mock.Anything, invalidAsset).
		Return(nil, fmt.Errorf("invalid asset data")).Once()

	asset, err := app.CreateAsset(context.Background(), invalidAsset)

	assert.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, appError.ErrCreatingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_GetAssetByID_Success(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	expectedAsset := &model.Asset{
		ID: "asset123",
	}

	mockAssetsRepository.On("GetAssetByID", mock.Anything, "asset123").
		Return(expectedAsset, nil).Once()

	asset, err := app.GetAssetByID(context.Background(), "asset123")

	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, expectedAsset.ID, asset.ID)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_GetAssetByID_Failure(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	mockAssetsRepository.On("GetAssetByID", mock.Anything, "asset123").
		Return(nil, fmt.Errorf("database error")).Once()

	asset, err := app.GetAssetByID(context.Background(), "asset123")

	assert.Error(t, err)
	assert.Nil(t, asset)
	assert.Equal(t, appError.ErrGettingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_GetAssets_Success(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	expectedAssets := []*model.Asset{
		{ID: "asset123"},
		{ID: "asset456"},
	}

	mockAssetsRepository.On("GetAssets", mock.Anything, mock.AnythingOfType("*model.GetAssetsInput")).
		Return(expectedAssets, nil).Once()

	assets, err := app.GetAssets(context.Background(), &model.GetAssetsInput{})

	assert.NoError(t, err)
	assert.NotNil(t, assets)
	assert.Len(t, assets, 2)
	assert.Equal(t, expectedAssets[0].ID, assets[0].ID)
	assert.Equal(t, expectedAssets[1].ID, assets[1].ID)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_GetAssets_Failure(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	mockAssetsRepository.On("GetAssets", mock.Anything, mock.AnythingOfType("*model.GetAssetsInput")).
		Return(nil, fmt.Errorf("database error")).Once()

	assets, err := app.GetAssets(context.Background(), &model.GetAssetsInput{})

	assert.Error(t, err)
	assert.Nil(t, assets)
	assert.Equal(t, appError.ErrGettingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_UpdateAsset_Success(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetToUpdate := &model.Asset{
		ID: "asset123",
	}

	mockAssetsRepository.On("UpdateAsset", mock.Anything, assetToUpdate).
		Return(nil).Once()

	err := app.UpdateAsset(context.Background(), assetToUpdate)

	assert.NoError(t, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_UpdateAsset_AssetDoesNotBelongToUser(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetToUpdate := &model.Asset{
		ID: "asset123",
	}

	mockAssetsRepository.On("UpdateAsset", mock.Anything, assetToUpdate).
		Return(fmt.Errorf("does not exist")).Once()

	err := app.UpdateAsset(context.Background(), assetToUpdate)

	assert.Error(t, err)
	assert.Equal(t, appError.ErrAssetDoesNotBelongToTheUser, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_UpdateAsset_Failure(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetToUpdate := &model.Asset{
		ID: "asset123",
	}

	mockAssetsRepository.On("UpdateAsset", mock.Anything, assetToUpdate).
		Return(fmt.Errorf("general update error")).Once()

	err := app.UpdateAsset(context.Background(), assetToUpdate)

	assert.Error(t, err)
	assert.Equal(t, appError.ErrUpdatingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_DeleteAsset_Success(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetID := "asset123"
	address := "userAddress"

	mockAssetsRepository.On("DeleteAsset", mock.Anything, assetID, address).
		Return(nil).Once()

	err := app.DeleteAsset(context.Background(), assetID, address)

	assert.NoError(t, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_DeleteAsset_AssetDoesNotBelongToUser(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetID := "asset123"
	address := "userAddress"

	mockAssetsRepository.On("DeleteAsset", mock.Anything, assetID, address).
		Return(fmt.Errorf("does not exist")).Once()

	err := app.DeleteAsset(context.Background(), assetID, address)

	assert.Error(t, err)
	assert.Equal(t, appError.ErrAssetDoesNotBelongToTheUser, err)

	mockAssetsRepository.AssertExpectations(t)
}
func TestAssetsApp_DeleteAsset_Failure(t *testing.T) {
	mockAssetsRepository := new(assetsMock.Repository)
	mockLogger := logrus.New()

	app := app.NewAssetsApp(
		&config.Configuration{},
		nil,
		nil,
		mockAssetsRepository,
		nil,
		mockLogger,
	)

	assetID := "asset123"
	address := "userAddress"

	mockAssetsRepository.On("DeleteAsset", mock.Anything, assetID, address).
		Return(fmt.Errorf("general deletion error")).Once()

	err := app.DeleteAsset(context.Background(), assetID, address)

	assert.Error(t, err)
	assert.Equal(t, appError.ErrDeletingAsset, err)

	mockAssetsRepository.AssertExpectations(t)
}
