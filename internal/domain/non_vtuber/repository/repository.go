package repository

import "context"

// Repository contains functions for non-vtuber domain.
type Repository interface {
	Create(ctx context.Context, id int64) (int, error)
	GetAllIDs(ctx context.Context) ([]int64, int, error)
}
