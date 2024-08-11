package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/language/entity"
)

// Repository contains functions for language domain.
type Repository interface {
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	GetAll(ctx context.Context) ([]entity.Language, int, int, error)
	UpdateByID(ctx context.Context, id int64, data entity.Language) (int, error)
}
