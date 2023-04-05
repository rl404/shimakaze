package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/errors"
)

type agency struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	Member     int       `json:"member"`
	Subscriber int       `json:"subscriber"`
	UpdatedAt  time.Time `json:"updated_at"`
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
			ID:         a.ID,
			Name:       a.Name,
			Image:      a.Image,
			Member:     a.Member,
			Subscriber: a.Subscriber,
			UpdatedAt:  a.UpdatedAt,
		}
	}

	return res, http.StatusOK, nil
}

// GetAgencyByID to get agency by id.
func (s *service) GetAgencyByID(ctx context.Context, id int64) (*agency, int, error) {
	a, code, err := s.agency.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	return &agency{
		ID:         a.ID,
		Name:       a.Name,
		Image:      a.Image,
		Member:     a.Member,
		Subscriber: a.Subscriber,
		UpdatedAt:  a.UpdatedAt,
	}, http.StatusOK, nil
}
