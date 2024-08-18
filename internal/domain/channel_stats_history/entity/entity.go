package entity

import vtuberEntity "github.com/rl404/shimakaze/internal/domain/vtuber/entity"

// ChannelStats is entity for channel stats.
type ChannelStats struct {
	VtuberID    int64
	ChannelID   string
	ChannelType vtuberEntity.ChannelType
	Subscriber  int
}
