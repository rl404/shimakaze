package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors"
	_errors "github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get agency data.
// @tags Agency
// @produce json
// @param sort query string false "sort" enums(name,-name,member,-member,subscriber,-subscriber) default(name)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.agency}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /agencies [get]
func (api *API) handleGetAgencies(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	agencies, pagination, code, err := api.service.GetAgencies(r.Context(), service.GetAgenciesRequest{
		Sort:  sort,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, agencies, errors.Wrap(r.Context(), err), pagination)
}

// @summary Get agency data.
// @tags Agency
// @produce json
// @param id path integer true "wikia id"
// @success 200 {object} utils.Response{data=service.agency}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /agencies/{id} [get]
func (api *API) handleGetAgencyByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), err, _errors.ErrInvalidID))
		return
	}

	agency, code, err := api.service.GetAgencyByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, agency, errors.Wrap(r.Context(), err))
}
