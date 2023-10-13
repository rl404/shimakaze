package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get wikia image.
// @tags Wikia
// @produce json
// @param path path string true "wikia image url"
// @success 200 "PNG image"
// @success 400 {object} utils.Response
// @success 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /wikia/image/{path} [get]
func (api *API) handleGetWikiaImage(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	if path == "" {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), errors.ErrInvalidRequestFormat))
		return
	}

	image, code, err := api.service.GetWikiaImage(r.Context(), path)
	utils.ResponseWithPNG(w, code, image, stack.Wrap(r.Context(), err))
}
