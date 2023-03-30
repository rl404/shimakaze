package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/youtube/entity"
)

// Repository contains functions for youtube domain.
type Repository interface {
	GetChannelIDByURL(ctx context.Context, url string) (string, int, error)
	GetChannelByID(ctx context.Context, id string) (*entity.Channel, int, error)
	GetVideoIDsByChannelID(ctx context.Context, channelID string) ([]string, int, error)
	GetVideosByIDs(ctx context.Context, ids []string) ([]entity.Video, int, error)
}
