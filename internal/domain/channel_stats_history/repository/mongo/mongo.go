package mongo

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo contains functions for channel-stats-history mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new channel-stats-history mongodb.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("channel_stats_history"),
	}
}

// Create to create channel stats.
func (m *Mongo) Create(ctx context.Context, data entity.ChannelStats) (int, error) {
	if _, err := m.db.InsertOne(ctx, &channelStats{
		VtuberID:    data.VtuberID,
		ChannelID:   data.ChannelID,
		ChannelType: data.ChannelType,
		Subscriber:  data.Subscriber,
	}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusCreated, nil
}
