package service

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	agencyEntity "github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
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
	Subscriber          int             `json:"subscriber"`
	VideoCount          int             `json:"video_count"`
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
	UpdatedAt           time.Time       `json:"updated_at"`
}

type vtuberAgency struct {
	ID    int64  `json:"id" validate:"required,gte=1"`
	Name  string `json:"name" validate:"required" mod:"trim"`
	Image string `json:"image" validate:"url" mod:"trim"`
}

type vtuberChannel struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Type       entity.ChannelType `json:"type"`
	URL        string             `json:"url" validate:"required,url" mod:"trim"`
	Image      string             `json:"image"`
	Subscriber int                `json:"subscriber"`
	Videos     []vtuberVideo      `json:"videos"`
}

type vtuberVideo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Image     string     `json:"image"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

// GetVtuberByID to get vtuber by id.
func (s *service) GetVtuberByID(ctx context.Context, id int64) (*vtuber, int, error) {
	vt, code, err := s.vtuber.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
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
		videos := make([]vtuberVideo, len(c.Videos))
		for j, v := range c.Videos {
			videos[j] = vtuberVideo{
				ID:        v.ID,
				Title:     v.Title,
				URL:       v.URL,
				Image:     v.Image,
				StartDate: v.StartDate,
				EndDate:   v.EndDate,
			}
		}

		channels[i] = vtuberChannel{
			ID:         c.ID,
			Name:       c.Name,
			Type:       c.Type,
			URL:        c.URL,
			Image:      c.Image,
			Subscriber: c.Subscriber,
			Videos:     videos,
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
		Subscriber:          vt.Subscriber,
		VideoCount:          vt.VideoCount,
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
		UpdatedAt:           vt.UpdatedAt,
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
		return nil, code, stack.Wrap(ctx, err)
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(res), func(i, j int) {
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
		return nil, code, stack.Wrap(ctx, err)
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
		return nil, code, stack.Wrap(ctx, err)
	}

	agencies, _, code, err := s.agency.GetAll(ctx, agencyEntity.GetAllRequest{Page: 1, Limit: -1})
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
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
	Mode               entity.SearchMode    `validate:"oneof=all simple" mod:"default=all,trim,lcase"`
	Names              string               `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Name               string               `validate:"omitempty,gte=3" mod:"trim,lcase"`
	OriginalName       string               `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Nickname           string               `validate:"omitempty,gte=3" mod:"trim,lcase"`
	ExcludeActive      bool                 ``
	ExcludeRetired     bool                 ``
	StartDebutMonth    int                  `validate:"omitempty,gte=1"`
	EndDebutMonth      int                  `validate:"omitempty,gte=1"`
	StartDebutYear     int                  `validate:"omitempty,gte=1"`
	EndDebutYear       int                  `validate:"omitempty,gte=1"`
	StartRetiredMonth  int                  `validate:"omitempty,gte=1"`
	EndRetiredMonth    int                  `validate:"omitempty,gte=1"`
	StartRetiredYear   int                  `validate:"omitempty,gte=1"`
	EndRetiredYear     int                  `validate:"omitempty,gte=1"`
	Has2D              *bool                ``
	Has3D              *bool                ``
	CharacterDesigner  string               `mod:"trim"`
	Character2DModeler string               `mod:"trim"`
	Character3DModeler string               `mod:"trim"`
	InAgency           *bool                ``
	Agency             string               `mod:"trim"`
	AgencyID           int64                `validate:"omitempty,gte=1"`
	ChannelTypes       []entity.ChannelType `validate:"dive,gte=1" mod:"dive,trim"`
	BirthdayDay        int                  `validate:"omitempty,gte=1"`
	StartBirthdayMonth int                  `validate:"omitempty,gte=1"`
	EndBirthdayMonth   int                  `validate:"omitempty,gte=1"`
	BloodTypes         []string             `validate:"dive,gte=1" mod:"dive,trim"`
	Genders            []string             `validate:"dive,gte=1" mod:"dive,trim"`
	Zodiacs            []string             `validate:"dive,gte=1" mod:"dive,trim"`
	StartSubscriber    int                  `validate:"omitempty,gte=1"`
	EndSubscriber      int                  `validate:"omitempty,gte=1"`
	StartVideoCount    int                  `validate:"omitempty,gte=1"`
	EndVideoCount      int                  `validate:"omitempty,gte=1"`
	Sort               string               `validate:"oneof=name -name debut_date -debut_date retirement_date -retirement_date subscriber -subscriber video_count -video_count" mod:"default=name,trim,lcase"`
	Page               int                  `validate:"required,gte=1" mod:"default=1"`
	Limit              int                  `validate:"required,gte=-1" mod:"default=20"`
}

// GetVtubers to get vtuber list.
func (s *service) GetVtubers(ctx context.Context, data GetVtubersRequest) ([]vtuber, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	vtubers, total, code, err := s.vtuber.GetAll(ctx, entity.GetAllRequest{
		Mode:               data.Mode,
		Names:              data.Names,
		Name:               data.Name,
		OriginalName:       data.OriginalName,
		Nickname:           data.Nickname,
		ExcludeActive:      data.ExcludeActive,
		ExcludeRetired:     data.ExcludeRetired,
		StartDebutMonth:    data.StartDebutMonth,
		EndDebutMonth:      data.EndDebutMonth,
		StartDebutYear:     data.StartDebutYear,
		EndDebutYear:       data.EndDebutYear,
		StartRetiredMonth:  data.StartRetiredMonth,
		EndRetiredMonth:    data.EndRetiredMonth,
		StartRetiredYear:   data.StartRetiredYear,
		EndRetiredYear:     data.EndRetiredYear,
		Has2D:              data.Has2D,
		Has3D:              data.Has3D,
		CharacterDesigner:  data.CharacterDesigner,
		Character2DModeler: data.Character2DModeler,
		Character3DModeler: data.Character3DModeler,
		InAgency:           data.InAgency,
		Agency:             data.Agency,
		AgencyID:           data.AgencyID,
		ChannelTypes:       data.ChannelTypes,
		BirthdayDay:        data.BirthdayDay,
		StartBirthdayMonth: data.StartBirthdayMonth,
		EndBirthdayMonth:   data.EndBirthdayMonth,
		BloodTypes:         data.BloodTypes,
		Genders:            data.Genders,
		Zodiacs:            data.Zodiacs,
		StartSubscriber:    data.StartSubscriber,
		EndSubscriber:      data.EndSubscriber,
		StartVideoCount:    data.StartVideoCount,
		EndVideoCount:      data.EndVideoCount,
		Sort:               data.Sort,
		Page:               data.Page,
		Limit:              data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
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
			videos := make([]vtuberVideo, len(c.Videos))
			for j, v := range c.Videos {
				videos[j] = vtuberVideo{
					ID:        v.ID,
					Title:     v.Title,
					URL:       v.URL,
					Image:     v.Image,
					StartDate: v.StartDate,
					EndDate:   v.EndDate,
				}
			}

			channels[i] = vtuberChannel{
				ID:         c.ID,
				Name:       c.Name,
				Type:       c.Type,
				URL:        c.URL,
				Image:      c.Image,
				Subscriber: c.Subscriber,
				Videos:     videos,
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
			Subscriber:          vt.Subscriber,
			VideoCount:          vt.VideoCount,
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
			UpdatedAt:           vt.UpdatedAt,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}

// GetVtuberCharacterDesigners to get vtuber character designer list.
func (s *service) GetVtuberCharacterDesigners(ctx context.Context) ([]string, int, error) {
	designers, code, err := s.vtuber.GetCharacterDesigners(ctx)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return designers, http.StatusOK, nil
}

// GetVtuberCharacter2DModelers to get vtuber 2d modeler list.
func (s *service) GetVtuberCharacter2DModelers(ctx context.Context) ([]string, int, error) {
	modelers, code, err := s.vtuber.GetCharacter2DModelers(ctx)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return modelers, http.StatusOK, nil
}

// GetVtuberCharacter3DModelers to get vtuber 3d modeler list.
func (s *service) GetVtuberCharacter3DModelers(ctx context.Context) ([]string, int, error) {
	modelers, code, err := s.vtuber.GetCharacter3DModelers(ctx)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return modelers, http.StatusOK, nil
}
