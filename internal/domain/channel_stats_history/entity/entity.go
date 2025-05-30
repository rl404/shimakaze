package entity

import (
	"time"

	vtuberEntity "github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

// Group is history group.
type Group string

// Available history groups.
const (
	Daily   Group = "DAILY"
	Monthly Group = "MONTHLY"
	Yearly  Group = "YEARLY"
)

// ChannelStats is entity for channel stats.
type ChannelStats struct {
	VtuberID    int64
	ChannelID   string
	ChannelType vtuberEntity.ChannelType
	Subscriber  int
	CreatedAt   time.Time
}

// GetRequest is get request model.
type GetRequest struct {
	VtuberID  int64
	StartDate time.Time
	EndDate   time.Time
}
