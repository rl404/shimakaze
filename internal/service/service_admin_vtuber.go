package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/utils"
)

// VtuberOverriddenField is vtuber overridden field model.
type VtuberOverriddenField struct {
	ID             int64                  `json:"-" validate:"required,gte=1"`
	DebutDate      overriddenDate         `json:"debut_date"`
	RetirementDate overriddenDate         `json:"retirement_date"`
	Agencies       overriddenAgencies     `json:"agencies"`
	Affiliations   overriddenAffiliations `json:"affiliations"`
	Channels       overriddenChannels     `json:"channels"`
}

type overriddenDate struct {
	Flag  bool       `json:"flag"`
	Value *time.Time `json:"value"`
}

type overriddenAgencies struct {
	Flag  bool           `json:"flag"`
	Value []vtuberAgency `json:"value" validate:"dive" mod:"dive"`
}

type overriddenAffiliations struct {
	Flag  bool     `json:"flag"`
	Value []string `json:"value" validate:"dive,required" mod:"dive,trim"`
}

type overriddenChannels struct {
	Flag  bool            `json:"flag"`
	Value []vtuberChannel `json:"value" validate:"dive" mod:"dive"`
}

// DeleteVtuberByID to delete vtuber by id.
func (s *service) DeleteVtuberByID(ctx context.Context, id int64) (int, error) {
	code, err := s.vtuber.DeleteByID(ctx, id)
	if err != nil {
		return code, stack.Wrap(ctx, err)
	}

	if code, err := s.nonVtuber.Create(ctx, id); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

// ParseVtuberByID to request parse vtuber by id.
func (s *service) ParseVtuberByID(ctx context.Context, id int64) (int, error) {
	if err := s.publisher.PublishParseVtuber(ctx, id, true); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err)
	}
	return http.StatusAccepted, nil
}

// GetVtuberOverriddenFieldByID to get vtuber overridden field by id.
func (s *service) GetVtuberOverriddenFieldByID(ctx context.Context, id int64) (*VtuberOverriddenField, int, error) {
	vt, code, err := s.vtuber.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	agencies := make([]vtuberAgency, len(vt.OverriddenField.Agencies.Value))
	for i, a := range vt.OverriddenField.Agencies.Value {
		agencies[i] = vtuberAgency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]vtuberChannel, len(vt.OverriddenField.Channels.Value))
	for i, c := range vt.OverriddenField.Channels.Value {
		channels[i] = vtuberChannel{
			URL: c.URL,
		}
	}

	return &VtuberOverriddenField{
		DebutDate: overriddenDate{
			Flag:  vt.OverriddenField.DebutDate.Flag,
			Value: vt.OverriddenField.DebutDate.Value,
		},
		RetirementDate: overriddenDate{
			Flag:  vt.OverriddenField.RetirementDate.Flag,
			Value: vt.OverriddenField.RetirementDate.OldValue,
		},
		Agencies: overriddenAgencies{
			Flag:  vt.OverriddenField.Agencies.Flag,
			Value: agencies,
		},
		Affiliations: overriddenAffiliations{
			Flag:  vt.OverriddenField.Affiliations.Flag,
			Value: vt.OverriddenField.Affiliations.Value,
		},
		Channels: overriddenChannels{
			Flag:  vt.OverriddenField.Channels.Flag,
			Value: channels,
		},
	}, http.StatusOK, nil
}

// UpdateVtuberOverriddenFieldByID to update vtuber overridden field data by id.
func (s *service) UpdateVtuberOverriddenFieldByID(ctx context.Context, data VtuberOverriddenField) (int, error) {
	if err := utils.Validate(&data); err != nil {
		return http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	agencies := make([]entity.Agency, len(data.Agencies.Value))
	for i, a := range data.Agencies.Value {
		agencies[i] = entity.Agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]entity.Channel, len(data.Channels.Value))
	for i, c := range data.Channels.Value {
		channels[i] = entity.Channel{
			Type: entity.ParseChannelType(c.URL),
			URL:  c.URL,
		}
	}

	code, err := s.vtuber.UpdateOverriddenFieldByID(ctx, data.ID, entity.OverriddenField{
		DebutDate: entity.OverriddenDate{
			Flag:  data.DebutDate.Flag,
			Value: data.DebutDate.Value,
		},
		RetirementDate: entity.OverriddenDate{
			Flag:  data.RetirementDate.Flag,
			Value: data.RetirementDate.Value,
		},
		Agencies: entity.OverriddenAgencies{
			Flag:  data.Agencies.Flag,
			Value: agencies,
		},
		Affiliations: entity.OverriddenAffiliations{
			Flag:  data.Affiliations.Flag,
			Value: data.Affiliations.Value,
		},
		Channels: entity.OverriddenChannels{
			Flag:  data.Channels.Flag,
			Value: channels,
		},
	})
	if err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}
