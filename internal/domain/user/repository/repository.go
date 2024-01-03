package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/user/entity"
)

// Repository contains funtions for user domain.
type Repository interface {
	Upsert(ctx context.Context, data entity.User) (int, error)
	GetByID(ctx context.Context, id int64) (*entity.User, int, error)
}
