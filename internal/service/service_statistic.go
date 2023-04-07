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
