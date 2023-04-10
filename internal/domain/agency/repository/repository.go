package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/agency/entity"
)

// Repository contains functions for agency domain.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Agency, int, error)
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Agency, int, int, error)
	IsOld(ctx context.Context, id int64) (bool, int, error)
	UpdateByID(ctx context.Context, id int64, data entity.Agency) (int, error)
	GetOldIDs(ctx context.Context) ([]int64, int, error)
	GetCount(ctx context.Context) (int, int, error)
}
