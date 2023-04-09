package service

import (
	"context"
	"math"
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
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

type vtuberSubscriberCount struct {
	Min   int `json:"min"`
	Max   int `json:"max"`
	Count int `json:"count"`
}

// GetVtuberSubscriberCountRequest is get vtuber subscriber count request.
type GetVtuberSubscriberCountRequest struct {
	Interval int `validate:"required,gte=10000" mod:"default=100000"`
	Max      int `validate:"required,lte=5000000" mod:"default=5000000"`
}

// GetVtuberSubscriberCount to get vtuber subscriber count.
func (s *service) GetVtuberSubscriberCount(ctx context.Context, data GetVtuberSubscriberCountRequest) ([]vtuberSubscriberCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.GetSubscriberCount(ctx, data.Interval, data.Max)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberSubscriberCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberSubscriberCount{
			Min:   c.Min,
			Max:   c.Max,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuberDesignerCountRequest is get vtuber designer count request.
type GetVtuberDesignerCountRequest struct {
	Top int `validate:"required,gte=-1" mod:"default=10"`
}

type vtuberDesignerCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetVtuberDesignerCount to get vtuber character designer count.
func (s *service) GetVtuberDesignerCount(ctx context.Context, data GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.GetDesignerCount(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberDesignerCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberDesignerCount{
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuber2DModelerCount to get vtuber character 2d modeler count.
func (s *service) GetVtuber2DModelerCount(ctx context.Context, data GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.Get2DModelerCount(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberDesignerCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberDesignerCount{
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuber3DModelerCount to get vtuber character 3d modeler count.
func (s *service) GetVtuber3DModelerCount(ctx context.Context, data GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.Get3DModelerCount(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberDesignerCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberDesignerCount{
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuberAverageVideoCount to get vtuber average video count.
func (s *service) GetVtuberAverageVideoCount(ctx context.Context) (float64, int, error) {
	avg, code, err := s.vtuber.GetAverageVideoCount(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return avg, http.StatusOK, nil
}

// GetVtuberAverageVideoDuration to get vtuber average video duration.
func (s *service) GetVtuberAverageVideoDuration(ctx context.Context) (float64, int, error) {
	avg, code, err := s.vtuber.GetAverageVideoDuration(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return avg, http.StatusOK, nil
}

type vtuberVideoCount struct {
	Day   int `json:"day"` // 1=sunday 2=monday
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

// GetVtuberVideoCount to get vtuber video count.
func (s *service) GetVtuberVideoCount(ctx context.Context, hourly, daily bool) ([]vtuberVideoCount, int, error) {
	cnt, code, err := s.vtuber.GetVideoCount(ctx, hourly, daily)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberVideoCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberVideoCount{
			Day:   c.Day,
			Hour:  c.Hour,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}
