package entity

import "time"

// Agency is entity for agency.
type Agency struct {
	ID         int64
	Name       string
	Image      string
	Member     int
	Subscriber int
	UpdatedAt  time.Time
}

// GetAllRequest is entity for get all request.
type GetAllRequest struct {
	Sort  string
	Page  int
	Limit int
}
