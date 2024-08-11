package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
)

type language struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetLanguages to get language list.
func (s *service) GetLanguages(ctx context.Context) ([]language, *pagination, int, error) {
	languages, total, code, err := s.language.GetAll(ctx)
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]language, len(languages))
	for i, l := range languages {
		res[i] = language{
			ID:   l.ID,
			Name: l.Name,
		}
	}

	return res, &pagination{
		Page:  1,
		Limit: -1,
		Total: total,
	}, http.StatusOK, nil
}
