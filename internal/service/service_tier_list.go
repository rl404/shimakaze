package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/tier_list/entity"
	"github.com/rl404/shimakaze/internal/utils"
)

// TierList is tier list model.
type TierList struct {
	ID          string       `json:"id" mod:"trim"`
	Title       string       `json:"title" validate:"required,gte=3,lte=50" mod:"trim"`
	Description string       `json:"description" validate:"lte=1000" mod:"trim"`
	Tiers       []tier       `json:"tiers" validate:"required,dive" mod:"dive"`
	Options     []tierVtuber `json:"options" validate:"dive" mod:"dive"`
	User        User         `json:"user"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type tier struct {
	Label       string       `json:"label" validate:"required,lte=50" mod:"trim"`
	Description string       `json:"description" validate:"lte=1000" mod:"trim"`
	Color       string       `json:"color" validate:"required,lte=50" mod:"trim"`
	Size        string       `json:"size" validate:"required,lte=50" mod:"trim"`
	Vtubers     []tierVtuber `json:"vtubers" validate:"dive"`
}

type tierVtuber struct {
	ID          int64  `json:"id" validate:"required,gt=0"`
	Name        string `json:"name" validate:"required" mod:"trim"`
	Image       string `json:"image" mod:"trim"`
	Description string `json:"description" validate:"lte=1000" mod:"trim"`
}

func (s *service) tierListFromEntity(data entity.TierList) TierList {
	tiers := make([]tier, len(data.Tiers))
	for j, t := range data.Tiers {
		vtubers := make([]tierVtuber, len(t.Vtubers))
		for k, v := range t.Vtubers {
			vtubers[k] = tierVtuber{
				ID:          v.ID,
				Name:        v.Name,
				Image:       v.Image,
				Description: v.Description,
			}
		}

		tiers[j] = tier{
			Label:       t.Label,
			Description: t.Description,
			Color:       t.Color,
			Size:        t.Size,
			Vtubers:     vtubers,
		}
	}

	options := make([]tierVtuber, len(data.Options))
	for j, v := range data.Options {
		options[j] = tierVtuber{
			ID:          v.ID,
			Name:        v.Name,
			Image:       v.Image,
			Description: v.Description,
		}
	}

	return TierList{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Tiers:       tiers,
		Options:     options,
		User: User{
			ID:       data.User.ID,
			Username: data.User.Username,
		},
		UpdatedAt: data.UpdatedAt,
	}
}

// GetTierListsRequest is get tier lists request model.
type GetTierListsRequest struct {
	Query  string `validate:"omitempty,gte=3" mod:"trim,lcase"`
	UserID int64  `validate:"omitempty,gt=0"`
	Sort   string `validate:"oneof=title -title updated_at -updated_at" mod:"default=-updated_at,trim,lcase"`
	Page   int    `validate:"required,gte=1" mod:"default=1"`
	Limit  int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetTierLists to get tier lists.
func (s *service) GetTierLists(ctx context.Context, data GetTierListsRequest) ([]TierList, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	tierLists, total, code, err := s.tierList.Get(ctx, entity.GetRequest{
		Query:  data.Query,
		UserID: data.UserID,
		Sort:   data.Sort,
		Page:   data.Page,
		Limit:  data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]TierList, len(tierLists))
	for i, tl := range tierLists {
		res[i] = s.tierListFromEntity(tl)
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}

// GetTierListByID to get tier list by id.
func (s *service) GetTierListByID(ctx context.Context, id string) (*TierList, int, error) {
	tierList, code, err := s.tierList.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	res := s.tierListFromEntity(*tierList)
	return &res, http.StatusOK, nil
}

func (s *service) tierListToEntity(data TierList) entity.TierList {
	tiers := make([]entity.Tier, len(data.Tiers))
	for j, t := range data.Tiers {
		vtubers := make([]entity.Vtuber, len(t.Vtubers))
		for k, v := range t.Vtubers {
			vtubers[k] = entity.Vtuber{
				ID:          v.ID,
				Name:        v.Name,
				Image:       v.Image,
				Description: v.Description,
			}
		}

		tiers[j] = entity.Tier{
			Label:       t.Label,
			Description: t.Description,
			Color:       t.Color,
			Size:        t.Size,
			Vtubers:     vtubers,
		}
	}

	options := make([]entity.Vtuber, len(data.Options))
	for j, v := range data.Options {
		options[j] = entity.Vtuber{
			ID:          v.ID,
			Name:        v.Name,
			Image:       v.Image,
			Description: v.Description,
		}
	}

	return entity.TierList{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Tiers:       tiers,
		Options:     options,
		User: entity.User{
			ID:       data.User.ID,
			Username: data.User.Username,
		},
	}
}

// UpsertTierListByID to upsert tier list by id.
func (s *service) UpsertTierListByID(ctx context.Context, data TierList) (*TierList, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	tierList, code, err := s.tierList.UpsertByID(ctx, s.tierListToEntity(data))
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	res := s.tierListFromEntity(*tierList)
	return &res, code, nil
}

// DeleteTierListByID to delete tier list by id.
func (s *service) DeleteTierListByID(ctx context.Context, id string, userID int64) (int, error) {
	if code, err := s.tierList.DeleteByID(ctx, id, userID); err != nil {
		return code, stack.Wrap(ctx, err)
	}
	return http.StatusOK, nil
}
