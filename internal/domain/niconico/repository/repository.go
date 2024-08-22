package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/niconico/entity"
)

// Repository contains functions for niconico domain.
type Repository interface {
	GetUser(ctx context.Context, url string) (*entity.User, int, error)
	GetVideos(ctx context.Context, id string) ([]entity.Video, int, error)
	GetBroadcasts(ctx context.Context, id string) ([]entity.Video, int, error)
}
