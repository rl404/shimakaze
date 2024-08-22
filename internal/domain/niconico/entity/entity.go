package entity

import "time"

// User is entity for user.
type User struct {
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
	URL       string
}
