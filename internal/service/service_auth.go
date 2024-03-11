package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/token/entity"
	userEntity "github.com/rl404/shimakaze/internal/domain/user/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// AuthCallback is auth callback data.
type AuthCallback struct {
	Code string `json:"code" validate:"required"`
}

// Token is access & refresh token.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// HandleAuthCallback to handle auth callback.
func (s *service) HandleAuthCallback(ctx context.Context, data AuthCallback) (*Token, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	accessToken, code, err := s.sso.ExchangeCode(ctx, data.Code)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	user, code, err := s.sso.GetUser(ctx, accessToken)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	userDB := userEntity.User{
		ID:       user.ID,
		Username: user.Username,
	}

	if existingUser, _, _ := s.user.GetByID(ctx, user.ID); existingUser != nil {
		userDB = userEntity.User{
			ID:       existingUser.ID,
			Username: existingUser.Username,
			IsAdmin:  existingUser.IsAdmin,
		}
	}

	if code, err := s.user.Upsert(ctx, userDB); err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	token, code, err := s.generateToken(ctx, userDB.ID, userDB.Username, userDB.IsAdmin)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return token, http.StatusOK, nil
}

func (s *service) generateToken(ctx context.Context, userID int64, username string, isAdmin bool) (*Token, int, error) {
	accessUUID := utils.GenerateUUID()
	refreshUUID := utils.GenerateUUID()

	// Create access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, entity.CreateAccessTokenRequest{
		UserID:      userID,
		Username:    username,
		IsAdmin:     isAdmin,
		AccessUUID:  accessUUID,
		RefreshUUID: refreshUUID,
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	// Create refresh token.
	refreshToken, code, err := s.token.CreateRefreshToken(ctx, entity.CreateRefreshTokenRequest{
		UserID:      userID,
		Username:    username,
		IsAdmin:     isAdmin,
		RefreshUUID: refreshUUID,
	})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, http.StatusOK, nil
}

// ValidateToken to validate token.
func (s *service) ValidateToken(ctx context.Context, uuid string, userID int64) (int, error) {
	payload := s.token.GetToken(ctx, uuid)
	if payload.UserID != userID {
		return http.StatusUnauthorized, stack.Wrap(ctx, errors.ErrInvalidToken)
	}
	return http.StatusOK, nil
}

// InvalidateToken to invalidate token.
func (s *service) InvalidateToken(ctx context.Context, uuid string) (int, error) {
	payload := s.token.GetToken(ctx, uuid)

	// Delete access token.
	if code, err := s.token.DeleteToken(ctx, uuid); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Delete refresh token.
	if code, err := s.token.DeleteToken(ctx, payload.RefreshUUID); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

// JWTClaim is jwt claim.
type JWTClaim struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	IsAdmin     bool   `json:"is_admin"`
	AccessUUID  string `json:"-"`
	RefreshUUID string `json:"-"`
}

// RefreshToken to refresh token.
func (s *service) RefreshToken(ctx context.Context, data JWTClaim) (string, int, error) {
	// Create new access token.
	accessToken, code, err := s.token.CreateAccessToken(ctx, entity.CreateAccessTokenRequest{
		UserID:      data.UserID,
		Username:    data.Username,
		IsAdmin:     data.IsAdmin,
		AccessUUID:  utils.GenerateUUID(),
		RefreshUUID: data.RefreshUUID,
	})
	if err != nil {
		return "", code, stack.Wrap(ctx, err)
	}
	return accessToken, http.StatusOK, nil
}
