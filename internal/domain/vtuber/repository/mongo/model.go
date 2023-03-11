package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type vtuber struct {
	ID                  int64      `bson:"id"`
	Name                string     `bson:"name"`
	Image               string     `bson:"image"`
	OriginalNames       []string   `bson:"original_names"`
	Nicknames           []string   `bson:"nicknames"`
	Caption             string     `bson:"caption"`
	DebutDate           *time.Time `bson:"debut_date"`
	RetirementDate      *time.Time `bson:"retirement_date"`
	Has2D               bool       `bson:"has_2d"`
	Has3D               bool       `bson:"has_3d"`
	CharacterDesigners  []string   `bson:"character_designers"`
	Character2DModelers []string   `bson:"character_2d_modelers"`
	Character3DModelers []string   `bson:"character_3d_modelers"`
	Agencies            []agency   `bson:"agencies"`
	Affiliations        []string   `bson:"affiliations"`
	Channels            []channel  `bson:"channels"`
	SocialMedias        []string   `bson:"social_medias"`
	OfficialWebsites    []string   `bson:"official_websites"`
	Gender              string     `bson:"gender"`
	Age                 *float64   `bson:"age"`
	Birthday            *time.Time `bson:"birthday"`
	Height              *float64   `bson:"height"`
	Weight              *float64   `bson:"weight"`
	BloodType           string     `bson:"blood_type"`
	ZodiacSign          string     `bson:"zodiac_sign"`
	Emoji               string     `bson:"emoji"`
	CreatedAt           time.Time  `bson:"created_at"`
	UpdatedAt           time.Time  `bson:"updated_at"`
}

// MarshalBSON to override marshal function.
func (n *vtuber) MarshalBSON() ([]byte, error) {
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	n.UpdatedAt = time.Now()

	type n2 vtuber
	return bson.Marshal((*n2)(n))
}

func (v *vtuber) toEntity() *entity.Vtuber {
	agencies := make([]entity.Agency, len(v.Agencies))
	for i, a := range v.Agencies {
		agencies[i] = entity.Agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]entity.Channel, len(v.Channels))
	for i, c := range v.Channels {
		videos := make([]entity.Video, len(c.Videos))
		for j, vi := range c.Videos {
			videos[j] = entity.Video{
				Title:     vi.Title,
				URL:       vi.URL,
				Image:     vi.Image,
				StartDate: vi.StartDate,
				EndDate:   vi.EndDate,
			}
		}

		channels[i] = entity.Channel{
			Name:       c.Name,
			Type:       c.Type,
			URL:        c.URL,
			Image:      c.Image,
			Subscriber: c.Subscriber,
			Videos:     videos,
		}
	}

	return &entity.Vtuber{
		ID:                  v.ID,
		Name:                v.Name,
		Image:               v.Image,
		OriginalNames:       v.OriginalNames,
		Nicknames:           v.Nicknames,
		Caption:             v.Caption,
		DebutDate:           v.DebutDate,
		RetirementDate:      v.RetirementDate,
		Has2D:               v.Has2D,
		Has3D:               v.Has3D,
		CharacterDesigners:  v.CharacterDesigners,
		Character2DModelers: v.Character2DModelers,
		Character3DModelers: v.Character3DModelers,
		Agencies:            agencies,
		Affiliations:        v.Affiliations,
		Channels:            channels,
		SocialMedias:        v.SocialMedias,
		OfficialWebsites:    v.OfficialWebsites,
		Gender:              v.Gender,
		Age:                 v.Age,
		Birthday:            v.Birthday,
		Height:              v.Height,
		Weight:              v.Weight,
		BloodType:           v.BloodType,
		ZodiacSign:          v.ZodiacSign,
		Emoji:               v.Emoji,
	}
}

type agency struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type channel struct {
	Name       string             `bson:"name"`
	Type       entity.ChannelType `bson:"type"`
	URL        string             `bson:"url"`
	Image      string             `bson:"image"`
	Subscriber int                `bson:"subscriber"`
	Videos     []video            `bson:"videos"`
}

type video struct {
	Title     string    `bson:"title"`
	URL       string    `bson:"url"`
	Image     string    `bson:"image"`
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
}

func (m *Mongo) vtuberFromEntity(v entity.Vtuber) *vtuber {
	agencies := make([]agency, len(v.Agencies))
	for i, a := range v.Agencies {
		agencies[i] = agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]channel, len(v.Channels))
	for i, c := range v.Channels {
		videos := make([]video, len(c.Videos))
		for j, vid := range c.Videos {
			videos[j] = video{
				Title:     vid.Title,
				URL:       vid.URL,
				Image:     vid.Image,
				StartDate: vid.StartDate,
				EndDate:   vid.EndDate,
			}
		}

		channels[i] = channel{
			Name:       c.Name,
			Type:       c.Type,
			URL:        c.URL,
			Image:      c.Image,
			Subscriber: c.Subscriber,
			Videos:     videos,
		}
	}

	return &vtuber{
		ID:                  v.ID,
		Name:                v.Name,
		Image:               v.Image,
		OriginalNames:       v.OriginalNames,
		Nicknames:           v.Nicknames,
		Caption:             v.Caption,
		DebutDate:           v.DebutDate,
		RetirementDate:      v.RetirementDate,
		Has2D:               v.Has2D,
		Has3D:               v.Has3D,
		CharacterDesigners:  v.CharacterDesigners,
		Character2DModelers: v.Character2DModelers,
		Character3DModelers: v.Character3DModelers,
		Agencies:            agencies,
		Affiliations:        v.Affiliations,
		Channels:            channels,
		SocialMedias:        v.SocialMedias,
		OfficialWebsites:    v.OfficialWebsites,
		Gender:              v.Gender,
		Age:                 v.Age,
		Birthday:            v.Birthday,
		Height:              v.Height,
		Weight:              v.Weight,
		BloodType:           v.BloodType,
		ZodiacSign:          v.ZodiacSign,
		Emoji:               v.Emoji,
	}
}

func (m *Mongo) convertSort(sort string) bson.D {
	if sort == "" {
		return nil
	}

	if sort[0] == '-' {
		return bson.D{{Key: sort[1:], Value: -1}}
	}

	return bson.D{{Key: sort, Value: 1}}
}

func (m *Mongo) initFilter(filter bson.M, key string) {
	if filter[key] == nil {
		filter[key] = bson.M{}
	}
}

func (m *Mongo) getChannelTypeFilter(types []entity.ChannelType) bson.M {
	var includeTypes, excludeTypes []entity.ChannelType
	for _, t := range types {
		if t[0] == '-' {
			excludeTypes = append(excludeTypes, t[1:])
		} else {
			includeTypes = append(includeTypes, t)
		}
	}

	filter := bson.M{}
	if len(includeTypes) > 0 {
		filter["$all"] = includeTypes
	}

	if len(excludeTypes) > 0 {
		filter["$nin"] = excludeTypes
	}

	return filter
}
