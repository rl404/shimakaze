package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Delete vtuber data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param id path integer true "vtuber id"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/vtubers/{id} [delete]
func (api *API) handleDeleteVtuberByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	code, err := api.service.DeleteVtuberByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}

// @summary Parse vtuber data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param id path integer true "vtuber id"
// @success 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/vtubers/{id}/parse [post]
func (api *API) handleParseVtuberByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	code, err := api.service.ParseVtuberByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}

// @summary Get vtuber overridden field data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param id path integer true "vtuber id"
// @param request body service.VtuberOverriddenField true "request body"
// @success 200 {object} utils.Response{data=service.VtuberOverriddenField}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/vtubers/{id}/override [get]
func (api *API) handleGetVtuberOverriddenField(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	vtuber, code, err := api.service.GetVtuberOverriddenFieldByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, vtuber, stack.Wrap(r.Context(), err))
}

// @summary Update vtuber overridden field data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param id path integer true "vtuber id"
// @param request body service.VtuberOverriddenField true "request body"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/vtubers/{id}/override [put]
func (api *API) handleUpdateVtuberOverriddenField(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	var request service.VtuberOverriddenField
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidRequestData))
		return
	}

	request.ID = id

	code, err := api.service.UpdateVtuberOverriddenFieldByID(r.Context(), request)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}
