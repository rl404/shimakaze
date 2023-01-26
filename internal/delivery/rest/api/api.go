package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// API contains all functions for api endpoints.
type API struct {
	service service.Service
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(0), log.MiddlewareConfig{Error: true}))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(1), log.MiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Get("/wikia/image/*", api.maxConcurrent(api.handleGetWikiaImage, 5))

		r.Get("/vtubers/{id}", api.handleGetVtuberByID)
		r.Get("/vtubers/images", api.handleGetVtuberImages)
		r.Get("/vtubers/family-trees", api.handleGetVtuberFamilyTrees)
	})
}

// @summary Get wikia image.
// @tags Wikia
// @produce json
// @success 200 "PNG image"
// @success 400 {object} utils.Response
// @success 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /wikia/image/{path} [get]
func (api *API) handleGetWikiaImage(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	if path == "" {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidRequestFormat))
		return
	}

	image, code, err := api.service.GetWikiaImage(r.Context(), path)
	utils.ResponseWithPNG(w, code, image, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber data.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=service.vtuber}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /vtuber/{id} [get]
func (api *API) handleGetVtuberByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidID, err))
		return
	}

	vtuber, code, err := api.service.GetVtuberByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, vtuber, errors.Wrap(r.Context(), err))
}

// @summary Get all vtuber images.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=[]service.vtuberImage}
// @failure 500 {object} utils.Response
// @router /vtuber/images [get]
func (api *API) handleGetVtuberImages(w http.ResponseWriter, r *http.Request) {
	images, code, err := api.service.GetVtuberImages(r.Context())
	utils.ResponseWithJSON(w, code, images, errors.Wrap(r.Context(), err))
}

// @summary Get vtuber family trees.
// @tags Vtuber
// @produce json
// @success 200 {object} utils.Response{data=service.vtuberFamilyTree}
// @failure 500 {object} utils.Response
// @router /vtuber/family-trees [get]
func (api *API) handleGetVtuberFamilyTrees(w http.ResponseWriter, r *http.Request) {
	tree, code, err := api.service.GetVtuberFamilyTrees(r.Context())
	utils.ResponseWithJSON(w, code, tree, errors.Wrap(r.Context(), err))
}
