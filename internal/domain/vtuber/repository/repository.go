package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

// Repository contains functions for vtuber domain.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Vtuber, int, error)
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	GetAllImages(ctx context.Context, shuffle bool, limit int) ([]entity.Vtuber, int, error)
	UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error)
	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetOldIDs(ctx context.Context) ([]int64, int, error)
	GetAllForFamilyTree(ctx context.Context) ([]entity.Vtuber, int, error)
	GetAllForAgencyTree(ctx context.Context) ([]entity.Vtuber, int, error)
}
