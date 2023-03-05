package service

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

type vtuber struct {
	ID                  int64           `json:"id"`
	Name                string          `json:"name"`
	Image               string          `json:"image"`
	OriginalNames       []string        `json:"original_names"`
	Nicknames           []string        `json:"nicknames"`
	Caption             string          `json:"caption"`
	DebutDate           *time.Time      `json:"debut_date"`
	RetirementDate      *time.Time      `json:"retirement_date"`
	Has2D               bool            `json:"has_2d"`
	Has3D               bool            `json:"has_3d"`
	CharacterDesigners  []string        `json:"character_designers"`
	Character2DModelers []string        `json:"character_2d_modelers"`
	Character3DModelers []string        `json:"character_3d_modelers"`
	Agencies            []vtuberAgency  `json:"agencies"`
	Affiliations        []string        `json:"affiliations"`
	Channels            []vtuberChannel `json:"channels"`
	SocialMedias        []string        `json:"social_medias"`
	OfficialWebsites    []string        `json:"official_websites"`
	Gender              string          `json:"gender"`
	Age                 *float64        `json:"age"`
	Birthday            *time.Time      `json:"birthday"`
	Height              *float64        `json:"height"`
	Weight              *float64        `json:"weight"`
	BloodType           string          `json:"blood_type"`
	ZodiacSign          string          `json:"zodiac_sign"`
	Emoji               string          `json:"emoji"`
}

type vtuberAgency struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type vtuberChannel struct {
	Type entity.ChannelType `json:"type"`
	URL  string             `json:"url"`
}

// GetVtuberByID to get vtuber by id.
func (s *service) GetVtuberByID(ctx context.Context, id int64) (*vtuber, int, error) {
	vt, code, err := s.vtuber.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	agencies := make([]vtuberAgency, len(vt.Agencies))
	for i, a := range vt.Agencies {
		agencies[i] = vtuberAgency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]vtuberChannel, len(vt.Channels))
	for i, c := range vt.Channels {
		channels[i] = vtuberChannel{
			Type: c.Type,
			URL:  c.URL,
		}
	}

	return &vtuber{
		ID:                  vt.ID,
		Name:                vt.Name,
		Image:               vt.Image,
		OriginalNames:       vt.OriginalNames,
		Nicknames:           vt.Nicknames,
		Caption:             vt.Caption,
		DebutDate:           vt.DebutDate,
		RetirementDate:      vt.RetirementDate,
		Has2D:               vt.Has2D,
		Has3D:               vt.Has3D,
		CharacterDesigners:  vt.CharacterDesigners,
		Character2DModelers: vt.Character2DModelers,
		Character3DModelers: vt.Character3DModelers,
		Agencies:            agencies,
		Affiliations:        vt.Affiliations,
		Channels:            channels,
		SocialMedias:        vt.SocialMedias,
		OfficialWebsites:    vt.OfficialWebsites,
		Gender:              vt.Gender,
		Age:                 vt.Age,
		Birthday:            vt.Birthday,
		Height:              vt.Height,
		Weight:              vt.Weight,
		BloodType:           vt.BloodType,
		ZodiacSign:          vt.ZodiacSign,
		Emoji:               vt.Emoji,
	}, http.StatusOK, nil
}

type vtuberImage struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GetVtuberImages to get all vtuber images.
func (s *service) GetVtuberImages(ctx context.Context, shuffle bool, limit int) ([]vtuberImage, int, error) {
	images, code, err := s.vtuber.GetAllImages(ctx, shuffle, limit)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuberImage, len(images))
	for i, img := range images {
		res[i] = vtuberImage{
			ID:    img.ID,
			Name:  img.Name,
			Image: img.Image,
		}
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(res), func(i, j int) {
			res[i], res[j] = res[j], res[i]
		})
	}

	return res, http.StatusOK, nil
}

type vtuberFamilyTree struct {
	Nodes []vtuberFamilyTreeNode `json:"nodes"`
	Links []vtuberFamilyTreeLink `json:"links"`
}

type vtuberFamilyTreeNode struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	HasRetired bool   `json:"has_retired"`
}

type vtuberFamilyTreeLink struct {
	ID1  int64          `json:"id1"`
	ID2  int64          `json:"id2"`
	Role familyTreeRole `json:"role"`
}

type familyTreeRole string

const (
	familyTreeDesigner  familyTreeRole = "DESIGNER"
	familyTree2DModeler familyTreeRole = "2D_MODELER"
	familyTree3DModeler familyTreeRole = "3D_MODELER"
)

// GetVtuberFamilyTrees to get all vtuber family tree.
func (s *service) GetVtuberFamilyTrees(ctx context.Context) (*vtuberFamilyTree, int, error) {
	vtubers, code, err := s.vtuber.GetAllForFamilyTree(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	var tree vtuberFamilyTree

	parentIDMap := make(map[string]int64)

	for _, vtuber := range vtubers {
		tree.Nodes = append(tree.Nodes, vtuberFamilyTreeNode{
			ID:         vtuber.ID,
			Name:       vtuber.Name,
			Image:      vtuber.Image,
			HasRetired: vtuber.RetirementDate != nil,
		})

		for _, p := range vtuber.CharacterDesigners {
			if _, ok := parentIDMap[p]; !ok {
				parentIDMap[p] = -int64(len(parentIDMap) + 1)
			}

			tree.Links = append(tree.Links, vtuberFamilyTreeLink{
				ID1:  parentIDMap[p],
				ID2:  vtuber.ID,
				Role: familyTreeDesigner,
			})
		}

		for _, p := range vtuber.Character2DModelers {
			if _, ok := parentIDMap[p]; !ok {
				parentIDMap[p] = -int64(len(parentIDMap) + 1)
			}

			tree.Links = append(tree.Links, vtuberFamilyTreeLink{
				ID1:  parentIDMap[p],
				ID2:  vtuber.ID,
				Role: familyTree2DModeler,
			})
		}

		for _, p := range vtuber.Character3DModelers {
			if _, ok := parentIDMap[p]; !ok {
				parentIDMap[p] = -int64(len(parentIDMap) + 1)
			}

			tree.Links = append(tree.Links, vtuberFamilyTreeLink{
				ID1:  parentIDMap[p],
				ID2:  vtuber.ID,
				Role: familyTree3DModeler,
			})
		}
	}

	for k, v := range parentIDMap {
		tree.Nodes = append(tree.Nodes, vtuberFamilyTreeNode{
			ID:   v,
			Name: k,
		})
	}

	return &tree, http.StatusOK, nil
}

type vtuberAgencyTree struct {
	Nodes []vtuberAgencyTreeNode `json:"nodes"`
	Links []vtuberAgencyTreeLink `json:"links"`
}

type vtuberAgencyTreeNode struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	HasRetired bool     `json:"has_retired"`
	Agencies   []string `json:"agencies"`
}

type vtuberAgencyTreeLink struct {
	ID1 int64 `json:"id1"`
	ID2 int64 `json:"id2"`
}

// GetVtuberAgencyTrees to get all vtuber agency tree.
func (s *service) GetVtuberAgencyTrees(ctx context.Context) (*vtuberAgencyTree, int, error) {
	vtubers, code, err := s.vtuber.GetAllForAgencyTree(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	agencies, code, err := s.agency.GetAll(ctx)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	var tree vtuberAgencyTree

	for _, agency := range agencies {
		tree.Nodes = append(tree.Nodes, vtuberAgencyTreeNode{
			ID:    -agency.ID,
			Name:  agency.Name,
			Image: agency.Image,
		})
	}

	for _, vtuber := range vtubers {
		agenciesTmp := make([]string, len(vtuber.Agencies))
		for i, a := range vtuber.Agencies {
			agenciesTmp[i] = a.Name

			tree.Links = append(tree.Links, vtuberAgencyTreeLink{
				ID1: -a.ID,
				ID2: vtuber.ID,
			})
		}

		tree.Nodes = append(tree.Nodes, vtuberAgencyTreeNode{
			ID:         vtuber.ID,
			Name:       vtuber.Name,
			Image:      vtuber.Image,
			HasRetired: vtuber.RetirementDate != nil,
			Agencies:   agenciesTmp,
		})
	}

	return &tree, http.StatusOK, nil
}

// GetVtubersRequest is get vtubers request model.
type GetVtubersRequest struct {
	Mode  entity.SearchMode `validate:"oneof=all stats" mod:"default=all,trim,lcase"`
	Name  string            `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Page  int               `validate:"required,gte=1" mod:"default=1"`
	Limit int               `validate:"required,gte=-1" mod:"default=20"`
}

// GetVtubers to get vtuber list.
func (s *service) GetVtubers(ctx context.Context, data GetVtubersRequest) ([]vtuber, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	vtubers, total, code, err := s.vtuber.GetAll(ctx, entity.GetAllRequest{
		Mode:  data.Mode,
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]vtuber, len(vtubers))
	for i, vt := range vtubers {
		agencies := make([]vtuberAgency, len(vt.Agencies))
		for i, a := range vt.Agencies {
			agencies[i] = vtuberAgency{
				ID:    a.ID,
				Name:  a.Name,
				Image: a.Image,
			}
		}

		channels := make([]vtuberChannel, len(vt.Channels))
		for i, c := range vt.Channels {
			channels[i] = vtuberChannel{
				Type: c.Type,
				URL:  c.URL,
			}
		}

		res[i] = vtuber{
			ID:                  vt.ID,
			Name:                vt.Name,
			Image:               vt.Image,
			OriginalNames:       vt.OriginalNames,
			Nicknames:           vt.Nicknames,
			Caption:             vt.Caption,
			DebutDate:           vt.DebutDate,
			RetirementDate:      vt.RetirementDate,
			Has2D:               vt.Has2D,
			Has3D:               vt.Has3D,
			CharacterDesigners:  vt.CharacterDesigners,
			Character2DModelers: vt.Character2DModelers,
			Character3DModelers: vt.Character3DModelers,
			Agencies:            agencies,
			Affiliations:        vt.Affiliations,
			Channels:            channels,
			SocialMedias:        vt.SocialMedias,
			OfficialWebsites:    vt.OfficialWebsites,
			Gender:              vt.Gender,
			Age:                 vt.Age,
			Birthday:            vt.Birthday,
			Height:              vt.Height,
			Weight:              vt.Weight,
			BloodType:           vt.BloodType,
			ZodiacSign:          vt.ZodiacSign,
			Emoji:               vt.Emoji,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}
