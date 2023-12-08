package entity

import "time"

// TierList is entity for tier list.
type TierList struct {
	ID          string
	Title       string
	Description string
	Tiers       []Tier
	Options     []Vtuber
	User        User
	UpdatedAt   time.Time
}

// Tier is entity for tier.
type Tier struct {
	Label       string
	Description string
	Color       string
	Size        string
	Vtubers     []Vtuber
}

// Vtuber is entity for vtuber.
type Vtuber struct {
	ID          int64
	Name        string
	Image       string
	Description string
}

// User is entity for user.
type User struct {
	ID       int64
	Username string
}

// GetRequest is get request model.
type GetRequest struct {
	Query  string
	UserID int64
	Sort   string
	Page   int
	Limit  int
}
