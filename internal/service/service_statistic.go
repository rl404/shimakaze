package service

import (
	"context"
	"math"
	"net/http"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
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

type vtuberVideoCountByDate struct {
	Day   int `json:"day"` // 1=sunday 2=monday
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

// GetVtuberVideoCountByDate to get vtuber video count by date.
func (s *service) GetVtuberVideoCountByDate(ctx context.Context, hourly, daily bool) ([]vtuberVideoCountByDate, int, error) {
	cnt, code, err := s.vtuber.GetVideoCountByDate(ctx, hourly, daily)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberVideoCountByDate, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberVideoCountByDate{
			Day:   c.Day,
			Hour:  c.Hour,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberVideoCount struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetVtuberVideoCountRequest is get vtuber video count request.
type GetVtuberVideoCountRequest struct {
	Top int `validate:"required,gte=-1" mod:"default=10"`
}

// GetVtuberVideoCount to get vtuber video count.
func (s *service) GetVtuberVideoCount(ctx context.Context, data GetVtuberVideoCountRequest) ([]vtuberVideoCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.GetVideoCount(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberVideoCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberVideoCount{
			ID:    c.ID,
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberVideoDuration struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Duration float64 `json:"duration"` // second
}

// GetVtuberVideoDurationRequest is get vtuber video duration request.
type GetVtuberVideoDurationRequest struct {
	Top int `validate:"required,gte=-1" mod:"default=10"`
}

// GetVtuberVideoDuration to get vtuber video count.
func (s *service) GetVtuberVideoDuration(ctx context.Context, data GetVtuberVideoDurationRequest) ([]vtuberVideoDuration, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.GetVideoDuration(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberVideoDuration, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberVideoDuration{
			ID:       c.ID,
			Name:     c.Name,
			Duration: c.Duration,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberBirthdayCount struct {
	Month int `json:"month"`
	Day   int `json:"day"`
	Count int `json:"count"`
}

// GetVtuberBirthdayCount to get vtuber birthday count.
func (s *service) GetVtuberBirthdayCount(ctx context.Context) ([]vtuberBirthdayCount, int, error) {
	cnt, code, err := s.vtuber.GetBirthdayCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberBirthdayCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberBirthdayCount{
			Month: c.Month,
			Day:   c.Day,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetVtuberAverageHeight to get vtuber average height.
func (s *service) GetVtuberAverageHeight(ctx context.Context) (float64, int, error) {
	height, code, err := s.vtuber.GetAverageHeight(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return height, http.StatusOK, nil
}

// GetVtuberAverageWeight to get vtuber average weight.
func (s *service) GetVtuberAverageWeight(ctx context.Context) (float64, int, error) {
	weight, code, err := s.vtuber.GetAverageWeight(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}
	return weight, http.StatusOK, nil
}

// GetVtuberBloodTypeCountRequest is get vtuber blood type count request.
type GetVtuberBloodTypeCountRequest struct {
	Top int `validate:"required,gte=-1" mod:"default=5"`
}

type vtuberBloodTypeCount struct {
	BloodType string `json:"blood_type"`
	Count     int    `json:"count"`
}

func (s *service) GetVtuberBloodTypeCount(ctx context.Context, data GetVtuberBloodTypeCountRequest) ([]vtuberBloodTypeCount, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	cnt, code, err := s.vtuber.GetBloodTypeCount(ctx, data.Top)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberBloodTypeCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberBloodTypeCount{
			BloodType: c.BloodType,
			Count:     c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberChannelTypeCount struct {
	ChannelType entity.ChannelType `json:"channel_type"`
	Count       int                `json:"count"`
}

// GetVtuberChannelTypeCount to get vtuber channel type count.
func (s *service) GetVtuberChannelTypeCount(ctx context.Context) ([]vtuberChannelTypeCount, int, error) {
	cnt, code, err := s.vtuber.GetChannelTypeCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberChannelTypeCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberChannelTypeCount{
			ChannelType: c.ChannelType,
			Count:       c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberGenderCount struct {
	Gender string `json:"gender"`
	Count  int    `json:"count"`
}

// GetVtuberGenderCount to get vtuber gender count.
func (s *service) GetVtuberGenderCount(ctx context.Context) ([]vtuberGenderCount, int, error) {
	cnt, code, err := s.vtuber.GetGenderCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberGenderCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberGenderCount{
			Gender: c.Gender,
			Count:  c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type vtuberZodiacCount struct {
	Zodiac string `json:"zodiac"`
	Count  int    `json:"count"`
}

// GetVtuberZodiacCount to get vtuber zodiac count.
func (s *service) GetVtuberZodiacCount(ctx context.Context) ([]vtuberZodiacCount, int, error) {
	cnt, code, err := s.vtuber.GetZodiacCount(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberZodiacCount, len(cnt))
	for i, c := range cnt {
		res[i] = vtuberZodiacCount{
			Zodiac: c.Zodiac,
			Count:  c.Count,
		}
	}

	return res, http.StatusOK, nil
}
