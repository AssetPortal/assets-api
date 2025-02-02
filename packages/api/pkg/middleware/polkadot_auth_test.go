package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authMock "github.com/AssetPortal/assets-api/pkg/adapters/auth/mocks"
	tokensMock "github.com/AssetPortal/assets-api/pkg/adapters/tokens/mocks"
	"github.com/AssetPortal/assets-api/pkg/middleware"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPolkadotAuthMiddleware(t *testing.T) {
	tokensRepoMock := new(tokensMock.Repository)
	authClientMock := new(authMock.Client)

	headers := &model.AuthHeaders{
		Signature: "valid-signature",
		Address:   "valid-address",
		Message:   "valid-message",
	}

	polkadotAuth := middleware.NewPolkadotAuth(tokensRepoMock, authClientMock, true)

	tokensRepoMock.On("GetToken", mock.Anything, "valid-message").Return(&model.Token{Token: "valid-message"}, nil)

	authClientMock.On("VerifySignature", mock.Anything, "valid-message", "valid-address", "valid-signature").
		Return(&model.Auth{OK: true}, nil)

	tokensRepoMock.On("MarkTokenAsUsed", mock.Anything, "valid-message").Return(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, headers)
		assert.Equal(t, "valid-signature", headers.Signature)
		assert.Equal(t, "valid-address", headers.Address)
		assert.Equal(t, "valid-message", headers.Message)
		render.JSON(w, r, model.NewResponseError("Success"))
	})

	router := chi.NewRouter()
	router.With(
		httpin.NewInput(model.AuthHeaders{}),
	).With(polkadotAuth.Middleware).Get("/test", handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Address", "valid-address")
	req.Header.Set("X-Signature", "valid-signature")
	req.Header.Set("X-Message", "valid-message")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Success")

	tokensRepoMock.AssertExpectations(t)
	authClientMock.AssertExpectations(t)
}

func TestPolkadotAuthMiddlewareDisabled(t *testing.T) {
	tokensRepoMock := new(tokensMock.Repository)
	authClientMock := new(authMock.Client)
	polkadotAuth := middleware.NewPolkadotAuth(tokensRepoMock, authClientMock, false)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, model.NewResponseError("Success"))
	})

	router := chi.NewRouter()
	router.With(
		httpin.NewInput(model.AuthHeaders{}),
	).With(polkadotAuth.Middleware).Get("/test", handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Success")

	tokensRepoMock.AssertExpectations(t)
	authClientMock.AssertExpectations(t)
}
