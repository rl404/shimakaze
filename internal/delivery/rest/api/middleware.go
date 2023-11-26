package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

func (api *API) maxConcurrent(h http.HandlerFunc, n int) http.HandlerFunc {
	queue := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queue <- struct{}{}
		defer func() { <-queue }()

		h.ServeHTTP(w, r)
	})
}

type ctxJWTClaim struct{}
type ctxJWTToken struct{}

type tokenType int8

const (
	tokenAccess tokenType = iota + 1
	tokenRefresh
)

func (api *API) getJWTFromRequest(r *http.Request) string {
	// From query.
	query := r.URL.Query().Get("jwt")
	if query != "" {
		return query
	}

	// From header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	// From cookie.
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (api *API) jwtAuth(next http.HandlerFunc, tokenTypes ...tokenType) http.HandlerFunc {
	tokenType := tokenAccess
	if len(tokenTypes) > 0 {
		tokenType = tokenTypes[0]
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := api.getJWTFromRequest(r)
		if jwtToken == "" {
			utils.ResponseWithJSON(w, http.StatusUnauthorized, nil, stack.Wrap(r.Context(), errors.ErrInvalidToken))
			return
		}

		// Parse jwt.
		jwtClaim, code, err := api.parseJWT(r.Context(), jwtToken, tokenType)
		if err != nil {
			utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
			return
		}

		var uuid string
		switch tokenType {
		case tokenRefresh:
			uuid = jwtClaim.RefreshUUID
		default:
			uuid = jwtClaim.AccessUUID
		}

		// Validate token.
		if code, err := api.service.ValidateToken(r.Context(), uuid, jwtClaim.UserID); err != nil {
			utils.ResponseWithJSON(w, code, nil, stack.Wrap(r.Context(), err))
			return
		}

		ctx := context.WithValue(r.Context(), ctxJWTClaim{}, jwtClaim)
		ctx = context.WithValue(ctx, ctxJWTToken{}, jwtToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) parseJWT(ctx context.Context, jwtTokenStr string, tokenType tokenType) (*service.JWTClaim, int, error) {
	token, err := jwt.Parse(jwtTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInternalServer
		}
		switch tokenType {
		case tokenAccess:
			return []byte(api.accessSecret), nil
		case tokenRefresh:
			return []byte(api.refreshSecret), nil
		default:
			return nil, nil
		}
	})
	if err != nil {
		return nil, http.StatusUnauthorized, stack.Wrap(ctx, errors.ErrInvalidToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, errors.ErrInternalServer)
	}

	if !token.Valid {
		return nil, http.StatusUnauthorized, stack.Wrap(ctx, errors.ErrInvalidToken)
	}

	var t service.JWTClaim
	userID, _ := claims["sub"].(float64)
	t.UserID = int64(userID)
	t.AccessUUID, _ = claims["access_uuid"].(string)
	t.RefreshUUID, _ = claims["refresh_uuid"].(string)

	return &t, http.StatusOK, nil
}

func (api *API) getJWTClaimFromContext(ctx context.Context) (*service.JWTClaim, int, error) {
	claims, ok := ctx.Value(ctxJWTClaim{}).(*service.JWTClaim)
	if !ok {
		return nil, http.StatusInternalServerError, errors.ErrInternalServer
	}
	return claims, http.StatusOK, nil
}

func (api *API) getJWTTokenFromContext(ctx context.Context) (string, int, error) {
	token, ok := ctx.Value(ctxJWTToken{}).(string)
	if !ok {
		return "", http.StatusInternalServerError, errors.ErrInternalServer
	}
	return token, http.StatusOK, nil
}
