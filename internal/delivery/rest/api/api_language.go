package api

import (
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get language data.
// @tags Language
// @produce json
// @success 200 {object} utils.Response{data=[]service.language}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /languages [get]
func (api *API) handleGetLanguages(w http.ResponseWriter, r *http.Request) {
	languages, pagination, code, err := api.service.GetLanguages(r.Context())
	utils.ResponseWithJSON(w, code, languages, stack.Wrap(r.Context(), err), pagination)
}
