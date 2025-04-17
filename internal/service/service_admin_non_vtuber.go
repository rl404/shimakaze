package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/non_vtuber/entity"
	"github.com/rl404/shimakaze/internal/utils"
)

type nonVtuber struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetNonVtubersRequest is get non-vtuber list request model.
type GetNonVtubersRequest struct {
	Name  string `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetNonVtubers
func (s *service) GetNonVtubers(ctx context.Context, data GetNonVtubersRequest) ([]nonVtuber, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	nonVtubers, total, code, err := s.nonVtuber.GetAll(ctx, entity.GetAllRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]nonVtuber, len(nonVtubers))
	for i, v := range nonVtubers {
		res[i] = nonVtuber{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}

// DeleteNonVtuberByID
func (s *service) DeleteNonVtuberByID(ctx context.Context, id int64) (int, error) {
	if code, err := s.nonVtuber.DeleteByID(ctx, id); err != nil {
		return code, stack.Wrap(ctx, err)
	}
	return http.StatusOK, nil
}
