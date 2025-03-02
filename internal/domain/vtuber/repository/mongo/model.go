package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type vtuber struct {
	ID                  int64           `bson:"id"`
	Name                string          `bson:"name"`
	Image               string          `bson:"image"`
	OriginalNames       []string        `bson:"original_names"`
	Nicknames           []string        `bson:"nicknames"`
	Caption             string          `bson:"caption"`
	DebutDate           *time.Time      `bson:"debut_date"`
	RetirementDate      *time.Time      `bson:"retirement_date"`
	Has2D               bool            `bson:"has_2d"`
	Has3D               bool            `bson:"has_3d"`
	CharacterDesigners  []string        `bson:"character_designers"`
	Character2DModelers []string        `bson:"character_2d_modelers"`
	Character3DModelers []string        `bson:"character_3d_modelers"`
	Agencies            []agency        `bson:"agencies"`
	Affiliations        []string        `bson:"affiliations"`
	Languages           []language      `bson:"languages"`
	Channels            []channel       `bson:"channels"`
	Subscriber          int             `bson:"subscriber"`
	MonthlySubscriber   int             `bson:"monthly_subscriber"`
	VideoCount          int             `bson:"video_count"`
	AverageVideoLength  int             `bson:"average_video_length"`
	TotalVideoLength    int             `bson:"total_video_length"`
	SocialMedias        []string        `bson:"social_medias"`
	OfficialWebsites    []string        `bson:"official_websites"`
	Gender              string          `bson:"gender"`
	Age                 *float64        `bson:"age"`
	Birthday            *time.Time      `bson:"birthday"`
	Height              *float64        `bson:"height"`
	Weight              *float64        `bson:"weight"`
	BloodType           string          `bson:"blood_type"`
	ZodiacSign          string          `bson:"zodiac_sign"`
	Emoji               string          `bson:"emoji"`
	OverriddenField     overriddenField `bson:"overridden_field"`
	CreatedAt           time.Time       `bson:"created_at"`
	UpdatedAt           time.Time       `bson:"updated_at"`
}

type agency struct {
	ID    int64  `bson:"id"`
	Name  string `bson:"name"`
	Image string `bson:"image"`
}

type language struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
}

type channel struct {
	ID         string             `bson:"id"`
	Name       string             `bson:"name"`
	Type       entity.ChannelType `bson:"type"`
	URL        string             `bson:"url"`
	Image      string             `bson:"image"`
	Subscriber int                `bson:"subscriber"`
	Videos     []video            `bson:"videos"`
}

type video struct {
	ID        string     `bson:"id"`
	Title     string     `bson:"title"`
	URL       string     `bson:"url"`
	Image     string     `bson:"image"`
	StartDate *time.Time `bson:"start_date"`
	EndDate   *time.Time `bson:"end_date"`
}

type vtuberVideo struct {
	VtuberID       int64              `bson:"vtuber_id"`
	VtuberName     string             `bson:"vtuber_name"`
	VtuberImage    string             `bson:"vtuber_image"`
	ChannelID      string             `bson:"channel_id"`
	ChannelName    string             `bson:"channel_name"`
	ChannelType    entity.ChannelType `bson:"channel_type"`
	ChannelURL     string             `bson:"channel_url"`
	VideoID        string             `bson:"video_id"`
	VideoTitle     string             `bson:"video_title"`
	VideoURL       string             `bson:"video_url"`
	VideoImage     string             `bson:"video_image"`
	VideoStartDate *time.Time         `bson:"video_start_date"`
	VideoEndDate   *time.Time         `bson:"video_end_date"`
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

	languages := make([]entity.Language, len(v.Languages))
	for i, l := range v.Languages {
		languages[i] = entity.Language{
			ID:   l.ID,
			Name: l.Name,
		}
	}

	channels := make([]entity.Channel, len(v.Channels))
	for i, c := range v.Channels {
		videos := make([]entity.Video, len(c.Videos))
		for j, vi := range c.Videos {
			videos[j] = entity.Video{
				ID:        vi.ID,
				Title:     vi.Title,
				URL:       vi.URL,
				Image:     vi.Image,
				StartDate: vi.StartDate,
				EndDate:   vi.EndDate,
			}
		}

		channels[i] = entity.Channel{
			ID:         c.ID,
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
		Languages:           languages,
		Channels:            channels,
		Subscriber:          v.Subscriber,
		MonthlySubscriber:   v.MonthlySubscriber,
		VideoCount:          v.VideoCount,
		AverageVideoLength:  v.AverageVideoLength,
		TotalVideoLength:    v.TotalVideoLength,
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
		OverriddenField:     v.OverriddenField.toEntity(),
		UpdatedAt:           v.UpdatedAt,
	}
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

	languages := make([]language, len(v.Languages))
	for i, l := range v.Languages {
		languages[i] = language{
			ID:   l.ID,
			Name: l.Name,
		}
	}

	channels := make([]channel, len(v.Channels))
	for i, c := range v.Channels {
		videos := make([]video, len(c.Videos))
		for j, vid := range c.Videos {
			videos[j] = video{
				ID:        vid.ID,
				Title:     vid.Title,
				URL:       vid.URL,
				Image:     vid.Image,
				StartDate: vid.StartDate,
				EndDate:   vid.EndDate,
			}
		}

		channels[i] = channel{
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
		Languages:           languages,
		Channels:            channels,
		Subscriber:          v.Subscriber,
		MonthlySubscriber:   v.MonthlySubscriber,
		VideoCount:          v.VideoCount,
		AverageVideoLength:  v.AverageVideoLength,
		TotalVideoLength:    v.TotalVideoLength,
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
		OverriddenField:     m.overiddenFieldFromEntity(v.OverriddenField),
	}
}

func (m *Mongo) convertSort(sort string) bson.D {
	if sort == "" {
		sort = "name"
	}

	if sort[0] == '-' {
		if sort[1:] == "video_count" {
			return bson.D{{Key: sort[1:], Value: -1}, {Key: "retirement_date", Value: 1}, {Key: "id", Value: 1}}
		}
		return bson.D{{Key: sort[1:], Value: -1}, {Key: "id", Value: 1}}
	}

	if sort == "video_count" {
		return bson.D{{Key: sort, Value: 1}, {Key: "retirement_date", Value: -1}, {Key: "id", Value: 1}}
	}

	if sort == "debut_date" {
		return bson.D{{Key: "is_debut_date_null", Value: 1}, {Key: sort, Value: 1}, {Key: "id", Value: 1}}
	}

	return bson.D{{Key: sort, Value: 1}, {Key: "id", Value: 1}}
}

func (m *Mongo) getChannelTypeFilter(types []entity.ChannelType) bson.M {
	values := make([]string, len(types))
	for i, t := range types {
		values[i] = string(t)
	}
	return m.getArrayFilter(values)
}

func (m *Mongo) getArrayFilter(values []string) bson.M {
	var includeValues, excludeValues []string
	for _, v := range values {
		if v[0] == '-' {
			excludeValues = append(excludeValues, v[1:])
		} else {
			includeValues = append(includeValues, v)
		}
	}

	filter := bson.M{}
	if len(includeValues) > 0 {
		filter["$all"] = includeValues
	}

	if len(excludeValues) > 0 {
		filter["$nin"] = excludeValues
	}

	return filter
}

func (m *Mongo) getPipeline(stages ...bson.D) mongo.Pipeline {
	var pipelines mongo.Pipeline
	for _, stage := range stages {
		if len(stage) > 0 {
			pipelines = append(pipelines, stage)
		}
	}
	return pipelines
}

func (m *Mongo) addStage(stageKey string, stages bson.D, key string, value interface{}) bson.D {
	for i, stage := range stages {
		if stage.Key != stageKey {
			continue
		}

		matchValue, ok := stage.Value.(bson.M)
		if !ok {
			continue
		}

		if matchValue[key] == nil {
			matchValue[key] = bson.M{}
		}

		if mValue, ok := value.(bson.M); ok {
			for k, v := range mValue {
				matchValue[key].(bson.M)[k] = v
			}
		} else {
			matchValue[key] = value
		}

		stages[i].Value = matchValue
		return stages
	}

	return append(stages, bson.E{
		Key:   stageKey,
		Value: bson.M{key: value},
	})
}

func (m *Mongo) addMatch(matchStage bson.D, key string, value interface{}) bson.D {
	return m.addStage("$match", matchStage, key, value)
}

func (m *Mongo) addField(projStage bson.D, key string, value interface{}) bson.D {
	return m.addStage("$addFields", projStage, key, value)
}

func (m *Mongo) mergeDebutRetiredMonthly(debut, retired []statusCountMonthly) []entity.DebutRetireCount {
	debutMap := make(map[int]map[int]int)
	retiredMap := make(map[int]map[int]int)

	minMonth := 1
	maxMonth := 12
	minYear := time.Now().Year()
	maxYear := time.Now().Year()

	for _, d := range debut {
		if d.Year == 0 || d.Month == 0 {
			continue
		}

		if debutMap[d.Year] == nil {
			debutMap[d.Year] = make(map[int]int)
		}

		debutMap[d.Year][d.Month] = d.Count

		if d.Year < minYear {
			minYear = d.Year
		}
	}

	for _, r := range retired {
		if r.Year == 0 || r.Month == 0 {
			continue
		}

		if retiredMap[r.Year] == nil {
			retiredMap[r.Year] = make(map[int]int)
		}

		retiredMap[r.Year][r.Month] = r.Count

		if r.Year < minYear {
			minYear = r.Year
		}
	}

	var data []entity.DebutRetireCount
	for y := minYear; y <= maxYear; y++ {
		for m := minMonth; m <= maxMonth; m++ {
			data = append(data, entity.DebutRetireCount{
				Year:   y,
				Month:  m,
				Debut:  debutMap[y][m],
				Retire: retiredMap[y][m],
			})
		}
	}

	return data
}

func (m *Mongo) mergeDebutRetiredYearly(debut, retired []statusCountYearly) []entity.DebutRetireCount {
	debutMap := make(map[int]int)
	retiredMap := make(map[int]int)

	minYear := time.Now().Year()
	maxYear := time.Now().Year()

	for _, d := range debut {
		if d.Year == 0 {
			continue
		}

		debutMap[d.Year] = d.Count

		if d.Year < minYear {
			minYear = d.Year
		}
	}

	for _, r := range retired {
		if r.Year == 0 {
			continue
		}

		retiredMap[r.Year] = r.Count

		if r.Year < minYear {
			minYear = r.Year
		}
	}

	var data []entity.DebutRetireCount
	for y := minYear; y <= maxYear; y++ {
		data = append(data, entity.DebutRetireCount{
			Year:   y,
			Debut:  debutMap[y],
			Retire: retiredMap[y],
		})
	}

	return data
}
