package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/twitch/entity"
)

// Repository contains functions for twitch domain.
type Repository interface {
	GetUser(ctx context.Context, name string) (*entity.User, int, error)
	GetFollowerCount(ctx context.Context, id string) (int, int, error)
	GetVideos(ctx context.Context, id string) ([]entity.Video, int, error)
}
