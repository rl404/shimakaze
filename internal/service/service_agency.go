package service

import (
	"context"
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
)

type agency struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GetAgencies to get agency list.
func (s *service) GetAgencies(ctx context.Context) ([]agency, int, error) {
	agencies, code, err := s.agency.GetAll(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]agency, len(agencies))
	for i, a := range agencies {
		res[i] = agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	return res, http.StatusOK, nil
}
