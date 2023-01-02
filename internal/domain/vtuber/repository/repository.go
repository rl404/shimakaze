package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

// Repository contains functions for vtuber domain.
type Repository interface {
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error)
	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetOldIDs(ctx context.Context) ([]int64, int, error)
}
