package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_cache "github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/token/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

type cache struct {
	cacher         _cache.Cacher
	accessSecret   string
	accessExpired  time.Duration
	refreshSecret  string
	refreshExpired time.Duration
}

// New to create new token cache.
func New(cacher _cache.Cacher,
	as string, ae time.Duration,
	rs string, re time.Duration,
) *cache {
	return &cache{
		cacher:         cacher,
		accessSecret:   as,
		accessExpired:  ae,
		refreshSecret:  rs,
		refreshExpired: re,
	}
}

// CreateAccessToken to create new access token.
func (c *cache) CreateAccessToken(ctx context.Context, data entity.CreateAccessTokenRequest) (string, int, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaim := accessToken.Claims.(jwt.MapClaims)
	accessClaim["sub"] = data.UserID
	accessClaim["username"] = data.Username
	accessClaim["access_uuid"] = data.AccessUUID
	accessClaim["iat"] = time.Now().UTC().Unix()
	accessClaim["exp"] = time.Now().UTC().Add(c.accessExpired).Unix()
	accessTokenStr, err := accessToken.SignedString([]byte(c.accessSecret))
	if err != nil {
		return "", http.StatusInternalServerError, stack.Wrap(ctx, err)
	}

	keyAccess := utils.GetKey("token", data.AccessUUID)
	if err := c.cacher.Set(ctx, keyAccess, entity.Payload{UserID: data.UserID, RefreshUUID: data.RefreshUUID}, c.accessExpired); err != nil {
		return "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return accessTokenStr, http.StatusOK, nil
}

// CreateRefreshToken to create new refresh token.
func (c *cache) CreateRefreshToken(ctx context.Context, data entity.CreateRefreshTokenRequest) (string, int, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaim := refreshToken.Claims.(jwt.MapClaims)
	refreshClaim["sub"] = data.UserID
	refreshClaim["username"] = data.Username
	refreshClaim["refresh_uuid"] = data.RefreshUUID
	refreshClaim["iat"] = time.Now().UTC().Unix()
	refreshClaim["exp"] = time.Now().UTC().Add(c.refreshExpired).Unix()
	refreshTokenStr, err := refreshToken.SignedString([]byte(c.refreshSecret))
	if err != nil {
		return "", http.StatusInternalServerError, stack.Wrap(ctx, err)
	}

	keyRefresh := utils.GetKey("token", data.RefreshUUID)
	if err := c.cacher.Set(ctx, keyRefresh, entity.Payload{UserID: data.UserID}, c.refreshExpired); err != nil {
		return "", http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return refreshTokenStr, http.StatusOK, nil
}

// GetToken to get token from cache.
func (c *cache) GetToken(ctx context.Context, token string) (payload entity.Payload) {
	c.cacher.Get(ctx, utils.GetKey("token", token), &payload)
	return
}

// DeleteToken to delete token.
func (c *cache) DeleteToken(ctx context.Context, token string) (int, error) {
	if err := c.cacher.Delete(ctx, utils.GetKey("token", token)); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}
	return http.StatusOK, nil
}
