package service

import (
	"context"
	"math"
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
)

// GetVtuberCount to get vtuber count.
func (s *service) GetVtuberCount(ctx context.Context) (int, int, error) {
	cnt, code, err := s.vtuber.GetCount(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return cnt, http.StatusOK, nil
}

// GetAgencyCount to get agency count.
func (s *service) GetAgencyCount(ctx context.Context) (int, int, error) {
	cnt, code, err := s.agency.GetCount(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return cnt, http.StatusOK, nil
}

// GetVtuberAverageActiveTime to get vtuber average active time.
func (s *service) GetVtuberAverageActiveTime(ctx context.Context) (float64, int, error) {
	avg, code, err := s.vtuber.GetAverageActiveTime(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return math.Round(avg*100) / 100, http.StatusOK, nil
}

type vtuberStatusCount struct {
	Active  int `json:"active"`
	Retired int `json:"retired"`
}

// GetVtuberStatusCount to get vtuber status count.
func (s *service) GetVtuberStatusCount(ctx context.Context) (*vtuberStatusCount, int, error) {
	cnt, code, err := s.vtuber.GetStatusCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return &vtuberStatusCount{
		Active:  cnt.Active,
		Retired: cnt.Retired,
	}, http.StatusOK, nil
}

type vtuberDebutRetireCount struct {
	Year   int `json:"year"`
	Month  int `json:"month,omitempty"`
	Debut  int `json:"debut"`
	Retire int `json:"retire"`
}

// GetVtuberDebutRetireCountMonthly to get vtuber debut & retire count monthly.
func (s *service) GetVtuberDebutRetireCountMonthly(ctx context.Context) ([]vtuberDebutRetireCount, int, error) {
	cnt, code, err := s.vtuber.GetDebutRetireCountMonthly(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberDebutRetireCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberDebutRetireCount{
			Year:   c.Year,
			Month:  c.Month,
			Debut:  c.Debut,
			Retire: c.Retire,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuberDebutRetireCountYearly to get vtuber debut & retire count yearly.
func (s *service) GetVtuberDebutRetireCountYearly(ctx context.Context) ([]vtuberDebutRetireCount, int, error) {
	cnt, code, err := s.vtuber.GetDebutRetireCountYearly(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberDebutRetireCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberDebutRetireCount{
			Year:   c.Year,
			Debut:  c.Debut,
			Retire: c.Retire,
		}
	}

	return res, http.StatusOK, nil
}
