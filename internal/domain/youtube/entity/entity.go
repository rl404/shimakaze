package entity

import "time"

// Channel is entity for channel.
type Channel struct {
	ID         string
	Name       string
	Image      string
	Subscriber int
}

// Video is entity for video.
type Video struct {
	ID        string
	Title     string
	Image     string
	StartDate *time.Time
	EndDate   *time.Time
}
