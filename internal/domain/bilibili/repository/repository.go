package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/bilibili/entity"
)

// Repository contains functions for bilibili domain.
type Repository interface {
	GetUser(ctx context.Context, id string) (*entity.User, int, error)
	GetFollowerCount(ctx context.Context, id string) (int, int, error)
	GetVideos(ctx context.Context, id string) ([]entity.Video, int, error)
}
