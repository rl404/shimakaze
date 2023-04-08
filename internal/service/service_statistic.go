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

type vtuberModelCount struct {
	None      int `json:"none"`
	Has2DOnly int `json:"has_2d_only"`
	Has3DOnly int `json:"has_3d_only"`
	Both      int `json:"both"`
}

// GetVtuberModelCount to get vtuber 2d & 3d model count.
func (s *service) GetVtuberModelCount(ctx context.Context) (*vtuberModelCount, int, error) {
	cnt, code, err := s.vtuber.GetModelCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return &vtuberModelCount{
		None:      cnt.None,
		Has2DOnly: cnt.Has2DOnly,
		Has3DOnly: cnt.Has3DOnly,
		Both:      cnt.Both,
	}, http.StatusOK, nil
}

type vtuberInAgencyCount struct {
	InAgency    int `json:"in_agency"`
	NotInAgency int `json:"not_in_agency"`
}

// GetVtuberInAgencyCount to get vtuber in agency count.
func (s *service) GetVtuberInAgencyCount(ctx context.Context) (*vtuberInAgencyCount, int, error) {
	cnt, code, err := s.vtuber.GetInAgencyCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return &vtuberInAgencyCount{
		InAgency:    cnt.InAgency,
		NotInAgency: cnt.NotInAgency,
	}, http.StatusOK, nil
}
