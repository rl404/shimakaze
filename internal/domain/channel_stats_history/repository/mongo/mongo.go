package mongo

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Get to get channel stats.
func (m *Mongo) Get(ctx context.Context, data entity.GetRequest) ([]entity.ChannelStats, int, error) {
	match := bson.M{
		"vtuber_id": data.VtuberID,
		"$and": bson.A{
			bson.M{"created_at": bson.M{"$gte": primitive.NewDateTimeFromTime(data.StartDate)}},
			bson.M{"created_at": bson.M{"$lte": primitive.NewDateTimeFromTime(data.EndDate)}},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(0)

	c, err := m.db.Find(ctx, match, opts)
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	var histories []entity.ChannelStats
	for c.Next(ctx) {
		var history channelStats
		if err := c.Decode(&history); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		histories = append(histories, history.toEntity())
	}

	return histories, http.StatusOK, nil
}
