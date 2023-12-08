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

// @summary Get tier list data.
// @tags TierList
// @produce json
// @param query query string false "query"
// @param user_id query integer false "user id"
// @param sort query string false "sort" enums(title,-title,updated_at,-updated_at) default(-updated_at)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.TierList}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /tier-lists [get]
func (api *API) handleGetTierLists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	userID, _ := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	tierLists, pagination, code, err := api.service.GetTierLists(r.Context(), service.GetTierListsRequest{
		Query:  query,
		UserID: userID,
		Sort:   sort,
		Page:   page,
		Limit:  limit,
	})

	utils.ResponseWithJSON(w, code, tierLists, stack.Wrap(r.Context(), err), pagination)
}

// @summary Get tier list data.
// @tags TierList
// @produce json
// @param id path string true "tier list id"
// @success 200 {object} utils.Response{data=service.TierList}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /tier-lists/{id} [get]
func (api *API) handleGetTierListByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), errors.ErrInvalidID))
		return
	}

	tierList, code, err := api.service.GetTierListByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, tierList, stack.Wrap(r.Context(), err))
}

// @summary Upsert tier list data.
// @tags TierList
// @produce json
// @param request body service.TierList true "request body"
// @success 200 {object} utils.Response{data=service.TierList}
// @success 201 {object} utils.Response{data=service.TierList}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /tier-lists [put]
func (api *API) handleUpsertTierListByID(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	var request service.TierList
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidRequestData))
		return
	}

	request.User.ID = claims.UserID
	request.User.Username = claims.Username

	tierList, code, err := api.service.UpsertTierListByID(r.Context(), request)
	utils.ResponseWithJSON(w, code, tierList, stack.Wrap(r.Context(), err))
}

// @summary Delete tier list data.
// @tags TierList
// @produce json
// @param id path string true "tier list id"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 403 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /tier-lists/{id} [delete]
func (api *API) handleDeleteTierListByID(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), errors.ErrInvalidID))
		return
	}

	code, err = api.service.DeleteTierListByID(r.Context(), id, claims.UserID)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}
