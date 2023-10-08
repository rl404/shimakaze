package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/utils"
)

type agency struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	Member     int       `json:"member"`
	Subscriber int       `json:"subscriber"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetAgenciesRequest is get agencies request model.
type GetAgenciesRequest struct {
	Sort  string `validate:"oneof=name -name member -member subscriber -subscriber" mod:"default=name,trim,lcase"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetAgencies to get agency list.
func (s *service) GetAgencies(ctx context.Context, data GetAgenciesRequest) ([]agency, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	agencies, total, code, err := s.agency.GetAll(ctx, entity.GetAllRequest{
		Sort:  data.Sort,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
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

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
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
