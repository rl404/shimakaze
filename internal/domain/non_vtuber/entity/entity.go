package entity

// NonVtuber is entity for non-vtuber.
type NonVtuber struct {
	ID   int64
	Name string
}

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Name  string
	Page  int
	Limit int
}
