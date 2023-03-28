package entity

import "time"

// Agency is entity for agency.
type Agency struct {
	ID        int64
	Name      string
	Image     string
	Member    int
	UpdatedAt time.Time
}
