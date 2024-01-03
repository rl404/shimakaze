package api

import (
	"encoding/json"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// @summary Handle OAuth2 callback.
// @tags Auth
// @produce json
// @param request body service.AuthCallback true "request body"
// @success 200 {object} utils.Response{data=service.Token}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /auth/callback [post]
func (api *API) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	var request service.AuthCallback
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidRequestData))
		return
	}

	token, c, err := api.service.HandleAuthCallback(r.Context(), request)
	utils.ResponseWithJSON(w, c, token, stack.Wrap(r.Context(), err))
}

// @summary Refresh Token.
// @tags Auth
// @produce json
// @param Authorization header string true "Bearer jwt.refresh.token"
// @success 200 {object} utils.Response{data=service.Token}
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /auth/token/refresh [post]
func (api *API) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	token, code, err := api.getJWTTokenFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	accessToken, code, err := api.service.RefreshToken(r.Context(), *claims)
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	utils.ResponseWithJSON(w, code, service.Token{
		AccessToken:  accessToken,
		RefreshToken: token,
	}, stack.Wrap(r.Context(), err))
}

// @summary Logout.
// @tags Auth
// @produce json
// @param Authorization header string true "Bearer jwt.access.token"
// @success 200 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 401 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /auth/logout [post]
func (api *API) handleLogout(w http.ResponseWriter, r *http.Request) {
	claims, code, err := api.getJWTClaimFromContext(r.Context())
	if err != nil {
		utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
		return
	}

	code, err = api.service.InvalidateToken(r.Context(), claims.AccessUUID)
	utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
}
