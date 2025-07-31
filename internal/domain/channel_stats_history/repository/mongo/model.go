package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/channel_stats_history/entity"
	vtuberEntity "github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type channelStats struct {
	VtuberID    int64                    `bson:"vtuber_id"`
	ChannelID   string                   `bson:"channel_id"`
	ChannelType vtuberEntity.ChannelType `bson:"channel_type"`
	Subscriber  int                      `bson:"subscriber"`
	CreatedAt   time.Time                `bson:"created_at"`
}

// MarshalBSON to override marshal function.
func (cs *channelStats) MarshalBSON() ([]byte, error) {
	if cs.CreatedAt.IsZero() {
		cs.CreatedAt = time.Now()
	}

	type cs2 channelStats
	return bson.Marshal((*cs2)(cs))
}

func (cs *channelStats) toEntity() entity.ChannelStats {
	return entity.ChannelStats{
		VtuberID:    cs.VtuberID,
		ChannelID:   cs.ChannelID,
		ChannelType: cs.ChannelType,
		Subscriber:  cs.Subscriber,
		CreatedAt:   cs.CreatedAt,
	}
}
