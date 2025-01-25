package service

import (
	"fmt"
	"net/http"
	"time"

	appError "github.com/AssetPortal/assets-api/pkg/error"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
)

func (srv *Service) CreateToken(w http.ResponseWriter, r *http.Request) {
	token, err := srv.assetsApp.CreateToken(r.Context())
	if err != nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Error{
			Error: err.Error(),
		})
	} else {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, token)
	}
}

func (srv *Service) CreateAsset(w http.ResponseWriter, r *http.Request) {
	createAsset := r.Context().Value(httpin.Input).(*model.CreateAssetInput)
	if err := createAsset.Validate(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, model.NewResponseError(err.Error()))
		return
	}
	asset := &model.Asset{
		ID:          createAsset.ID,
		Description: createAsset.Description,
		Image:       createAsset.Image,
		Social:      createAsset.Social,
		Address:     createAsset.Address,
	}

	asset, err := srv.assetsApp.CreateAsset(r.Context(), asset)
	if err != nil {
		if err == appError.ErrCreatingAssetIDExists {
			render.Status(r, http.StatusUnprocessableEntity)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, model.NewResponseError(err.Error()))

	} else {
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, model.NewResponseData(asset))
	}
}

func (srv *Service) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	updateAsset := r.Context().Value(httpin.Input).(*model.UpdateAssetInput)
	if err := updateAsset.Validate(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, model.NewResponseError(err.Error()))
		return
	}
	asset := &model.Asset{
		ID:          updateAsset.ID,
		Description: updateAsset.Description,
		Image:       updateAsset.Image,
		Social:      updateAsset.Social,
		Address:     updateAsset.Address,
	}

	err := srv.assetsApp.UpdateAsset(r.Context(), asset)
	if err != nil {
		if err == appError.ErrCreatingAssetIDExists {
			render.Status(r, http.StatusUnprocessableEntity)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, model.NewResponseError(err.Error()))

	} else {
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, model.NewResponseData(asset))
	}
}

func (srv *Service) GetAssetByID(w http.ResponseWriter, r *http.Request) {
	getAssetByID := r.Context().Value(httpin.Input).(*model.GetAssetByIDInput)
	if err := getAssetByID.Validate(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, model.NewResponseError(err.Error()))
		return
	}
	asset, err := srv.assetsApp.GetAssetByID(r.Context(), getAssetByID.ID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.NewResponseError(err.Error()))
	} else if asset == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.NewResponseError("asset not found"))
	} else {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, model.NewResponseData(asset))
	}
}

func (srv *Service) GetAssets(w http.ResponseWriter, r *http.Request) {
	getAssets := r.Context().Value(httpin.Input).(*model.GetAssetsInput)
	if err := getAssets.Validate(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, model.NewResponseError(err.Error()))
		return
	}
	assets, err := srv.assetsApp.GetAssets(r.Context(), getAssets)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.NewResponseError(err.Error()))
	} else {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, model.NewResponseData(assets))
	}
}

func (srv *Service) UploadImage(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*model.UploadImageInput)
	if err := input.Validate(); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, model.NewResponseError(err.Error()))
		return
	}
	fileBytes, err := input.File.ReadAll()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.NewResponseError("failed to read file data"))
		return
	}
	contentType := http.DetectContentType(fileBytes)
	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}
	fileExt, valid := allowedTypes[contentType]
	if !valid {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.NewResponseError("invalid file type; only JPEG, PNG, and GIF are allowed"))
		return
	}
	fileKey := fmt.Sprintf("%s_%d%s", input.ID, time.Now().Unix(), fileExt)

	url, err := srv.assetsApp.UploadFile(r.Context(), fileKey, fileBytes, contentType)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.NewResponseError("failed to upload image"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, model.NewResponseData(url))
}
