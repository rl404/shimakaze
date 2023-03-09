package api

import (
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get agency data.
// @tags Agency
// @produce json
// @success 200 {object} utils.Response{data=[]service.agency}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /agencies [get]
func (api *API) handleGetAgencies(w http.ResponseWriter, r *http.Request) {
	agencies, code, err := api.service.GetAgencies(r.Context())
	utils.ResponseWithJSON(w, code, agencies, errors.Wrap(r.Context(), err))
}
