package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
)

// User is user model.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// GetProfile to get user profile.
func (s *service) GetProfile(ctx context.Context, userID int64) (*User, int, error) {
	user, code, err := s.user.GetByID(ctx, userID)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return &User{
		ID:       user.ID,
		Username: user.Username,
	}, http.StatusOK, nil
}
