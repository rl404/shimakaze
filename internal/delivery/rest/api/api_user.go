package api

import (
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Get user profile.
// @tags User
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @success 200 {object} utils.Response{data=service.User}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /profile [get]
func (api *API) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	user, code, err := api.service.GetProfile(r.Context(), claims.UserID)
	utils.ResponseWithJSON(w, code, user, stack.Wrap(r.Context(), err))
}
