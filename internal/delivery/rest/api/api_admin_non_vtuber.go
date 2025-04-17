package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get non-vtuber data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.nonVtuber}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/non-vtubers [get]
func (api *API) handleGetNonVtubers(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	nonVtubers, pagination, code, err := api.service.GetNonVtubers(r.Context(), service.GetNonVtubersRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, nonVtubers, stack.Wrap(r.Context(), err), pagination)
}

// @summary Delete non-vtuber data.
// @tags Admin
// @produce json
// @param Authorization header string true "Bearer jwt.admin_access.token"
// @param id path integer true "non-vtuber id"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /admin/non-vtubers/{id} [delete]
func (api *API) handleDeleteNonVtuberByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	code, err := api.service.DeleteNonVtuberByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}
