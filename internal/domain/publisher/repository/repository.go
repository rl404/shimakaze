package repository

import (
	"context"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseVtuber(ctx context.Context, id int64, forced bool) error
	PublishParseAgency(ctx context.Context, id int64, forced bool) error
}
