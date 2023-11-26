package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/sso/entity"
)

// Repository contains functions for sso domain.
type Repository interface {
	ExchangeCode(ctx context.Context, code string) (string, int, error)
	GetUser(ctx context.Context, token string) (*entity.User, int, error)
}
