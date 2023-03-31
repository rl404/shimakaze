package entity

import "time"

// User is entity for user.
type User struct {
	ID    string
	Name  string
	Image string
}

// Video is entity for video.
type Video struct {
	ID        string
	Title     string
	URL       string
	Image     string
	StartDate *time.Time
	EndDate   *time.Time
}
