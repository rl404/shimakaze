package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/tier_list/entity"
	"github.com/rl404/shimakaze/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tierList struct {
	ID          string    `bson:"id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Tiers       []tier    `bson:"tiers"`
	Options     []vtuber  `bson:"options"`
	CreatedBy   user      `bson:"created_by"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

type tier struct {
	Label       string   `bson:"label"`
	Description string   `bson:"description"`
	Color       string   `bson:"color"`
	Size        string   `bson:"size"`
	Vtubers     []vtuber `bson:"vtubers"`
}

type vtuber struct {
	ID          int64  `bson:"id"`
	Name        string `bson:"name"`
	Image       string `bson:"image"`
	Description string `bson:"description"`
}

type user struct {
	ID       int64  `bson:"id"`
	Username string `bson:"username"`
}

// MarshalBSON to override marshal function.
func (n *tierList) MarshalBSON() ([]byte, error) {
	if n.ID == "" {
		n.ID = utils.RandomStr(6)
	}

	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	n.UpdatedAt = time.Now()

	type n2 tierList
	return bson.Marshal((*n2)(n))
}

func (t *tierList) toEntity() *entity.TierList {
	tiers := make([]entity.Tier, len(t.Tiers))
	for i, tier := range t.Tiers {
		vtubers := make([]entity.Vtuber, len(tier.Vtubers))
		for j, v := range tier.Vtubers {
			vtubers[j] = entity.Vtuber{
				ID:          v.ID,
				Name:        v.Name,
				Image:       v.Image,
				Description: v.Description,
			}
		}

		tiers[i] = entity.Tier{
			Label:       tier.Label,
			Description: tier.Description,
			Color:       tier.Color,
			Size:        tier.Size,
			Vtubers:     vtubers,
		}
	}

	options := make([]entity.Vtuber, len(t.Options))
	for i, v := range t.Options {
		options[i] = entity.Vtuber{
			ID:          v.ID,
			Name:        v.Name,
			Image:       v.Image,
			Description: v.Description,
		}
	}

	return &entity.TierList{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Tiers:       tiers,
		Options:     options,
		User: entity.User{
			ID:       t.CreatedBy.ID,
			Username: t.CreatedBy.Username,
		},
		UpdatedAt: t.UpdatedAt,
	}
}

func (m *Mongo) fromEntity(data entity.TierList) *tierList {
	tiers := make([]tier, len(data.Tiers))
	for i, t := range data.Tiers {
		vtubers := make([]vtuber, len(t.Vtubers))
		for j, v := range t.Vtubers {
			vtubers[j] = vtuber{
				ID:          v.ID,
				Name:        v.Name,
				Image:       v.Image,
				Description: v.Description,
			}
		}

		tiers[i] = tier{
			Label:       t.Label,
			Description: t.Description,
			Color:       t.Color,
			Size:        t.Size,
			Vtubers:     vtubers,
		}
	}

	options := make([]vtuber, len(data.Options))
	for i, v := range data.Options {
		options[i] = vtuber{
			ID:          v.ID,
			Name:        v.Name,
			Image:       v.Image,
			Description: v.Description,
		}
	}

	return &tierList{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Tiers:       tiers,
		Options:     options,
		CreatedBy: user{
			ID:       data.User.ID,
			Username: data.User.Username,
		},
	}
}

func (m *Mongo) convertSort(sort string) bson.D {
	if sort == "" {
		sort = "-updated_at"
	}

	if sort[0] == '-' {
		return bson.D{{Key: sort[1:], Value: -1}}
	}

	return bson.D{{Key: sort, Value: 1}}
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
