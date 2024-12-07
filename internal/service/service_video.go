package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

type video struct {
	VtuberID       int64              `json:"vtuber_id"`
	VtuberName     string             `json:"vtuber_name"`
	VtuberImage    string             `json:"vtuber_image"`
	ChannelID      string             `json:"channel_id"`
	ChannelName    string             `json:"channel_name"`
	ChannelType    entity.ChannelType `json:"channel_type"`
	ChannelURL     string             `json:"channel_url"`
	VideoID        string             `json:"video_id"`
	VideoTitle     string             `json:"video_title"`
	VideoURL       string             `json:"video_url"`
	VideoImage     string             `json:"video_image"`
	VideoStartDate *time.Time         `json:"video_start_date"`
	VideoEndDate   *time.Time         `json:"video_end_date"`
}

// GetVideosRequest is get videos request model.
type GetVideosRequest struct {
	StartDate  string `validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate    string `validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	IsFinished *bool  ``
	Sort       string `validate:"oneof=video_start_date -video_start_date" mod:"default=-video_start_date,trim,lcase"`
	Page       int    `validate:"required,gte=1" mod:"default=1"`
	Limit      int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetVideos to get video list.
func (s *service) GetVideos(ctx context.Context, data GetVideosRequest) ([]video, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	var startDate *time.Time
	if data.StartDate != "" {
		tmp, err := time.Parse(time.RFC3339, data.StartDate)
		if err != nil {
			return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err, errors.ErrInvalidDate)
		}
		startDate = &tmp
	}

	var endDate *time.Time
	if data.EndDate != "" {
		tmp, err := time.Parse(time.RFC3339, data.EndDate)
		if err != nil {
			return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err, errors.ErrInvalidDate)
		}
		endDate = &tmp
	}

	videos, total, code, err := s.vtuber.GetVideos(ctx, entity.GetVideosRequest{
		StartDate:  startDate,
		EndDate:    endDate,
		IsFinished: data.IsFinished,
		Sort:       data.Sort,
		Page:       data.Page,
		Limit:      data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]video, len(videos))
	for i, v := range videos {
		res[i] = video{
			VtuberID:       v.VtuberID,
			VtuberName:     v.VtuberName,
			VtuberImage:    v.VtuberImage,
			ChannelID:      v.ChannelID,
			ChannelName:    v.ChannelName,
			ChannelType:    v.ChannelType,
			ChannelURL:     v.ChannelURL,
			VideoID:        v.VideoID,
			VideoTitle:     v.VideoTitle,
			VideoURL:       v.VideoURL,
			VideoImage:     v.VideoImage,
			VideoStartDate: v.VideoStartDate,
			VideoEndDate:   v.VideoEndDate,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}
