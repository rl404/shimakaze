package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/non_vtuber/entity"
)

// Repository contains functions for non-vtuber domain.
type Repository interface {
	Create(ctx context.Context, id int64, name string) (int, error)
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.NonVtuber, int, int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)
}
