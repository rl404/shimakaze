package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/token/entity"
)

// Repository contains functions for token domain.
type Repository interface {
	CreateAccessToken(ctx context.Context, data entity.CreateAccessTokenRequest) (string, int, error)
	CreateRefreshToken(ctx context.Context, data entity.CreateRefreshTokenRequest) (string, int, error)
	GetToken(ctx context.Context, token string) entity.Payload
	DeleteToken(ctx context.Context, token string) (int, error)
}
