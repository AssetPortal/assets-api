package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AssetPortal/assets-api/pkg/app"
	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/ggicci/httpin"
	httpin_core "github.com/ggicci/httpin/core"
	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/cors"

	polkadotMiddleware "github.com/AssetPortal/assets-api/pkg/middleware"
)

func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	if err2, ok := err.(*httpin_core.InvalidFieldError); ok {
		if err2.Key == "" {
			err = fmt.Errorf("invalid %s", err2.Directive)
		} else {
			err = fmt.Errorf("query argument \"%s\" has a wrong value", err2.Key)
		}
	} else {
		err = fmt.Errorf("error parsing query arguments")
	}
	render.Status(r, http.StatusUnprocessableEntity)
	render.JSON(rw, r, Error{
		Error: err.Error(),
	})
}
func init() {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)
	httpin_core.RegisterErrorHandler(CustomErrorHandler)
}

type Service struct {
	HTTPServer         *http.Server
	assetsApp          *app.AssetsApp
	polkadotMiddleware *polkadotMiddleware.PolkadotAuth
}

func NewService(assetsApp *app.AssetsApp, polkadotMiddleware *polkadotMiddleware.PolkadotAuth) *Service {
	return &Service{
		assetsApp:          assetsApp,
		polkadotMiddleware: polkadotMiddleware,
	}
}

func (srv *Service) Setup() {

	router := chi.NewRouter()
	cfg := srv.assetsApp.Config()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(cfg.HTTPTimeout) * time.Second))
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: fix
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		Debug:            false,
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		MaxAge:           300,
	}).Handler)

	router.Get("/nonce", srv.CreateToken)
	router.With(httpin.NewInput(model.UploadImageInput{})).Post("/upload", srv.UploadImage)
	router.With(
		httpin.NewInput(model.AuthHeaders{}),
	).With(srv.polkadotMiddleware.Middleware).With(
		httpin.NewInput(model.CreateAssetInput{}),
	).Post("/assets", srv.CreateAsset)
	router.With(
		httpin.NewInput(model.GetAssetsInput{}),
	).Get("/assets", srv.GetAssets)
	router.With(
		httpin.NewInput(model.GetAssetByIDInput{}),
	).Get("/assets/{id}", srv.GetAssetByID)
	// delete //TODO
	// update //TODO
	srv.HTTPServer = &http.Server{Addr: cfg.ServiceAddress, Handler: router}
}

func (srv *Service) Start() {
	fmt.Printf("Starting API server at %s\n", srv.assetsApp.Config().ServiceAddress)
	err := srv.HTTPServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// type responseWriter struct {
// 	http.ResponseWriter
// 	StatusCode int
// }

// func NewResponseWriter(w http.ResponseWriter) *responseWriter {
// 	return &responseWriter{w, http.StatusOK}
// }

// func (rw *responseWriter) WriteHeader(code int) {
// 	rw.StatusCode = code
// 	rw.ResponseWriter.WriteHeader(code)
// }
