package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/tier_list/entity"
)

// Repository contains functions for tier list domain.
type Repository interface {
	Get(ctx context.Context, data entity.GetRequest) ([]entity.TierList, int, int, error)
	GetByID(ctx context.Context, id string) (*entity.TierList, int, error)
	UpsertByID(ctx context.Context, data entity.TierList) (*entity.TierList, int, error)
	DeleteByID(ctx context.Context, id string, userID int64) (int, error)
}
