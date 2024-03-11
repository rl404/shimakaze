package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

type overriddenField struct {
	DebutDate      overriddenDate         `bson:"debut_date"`
	RetirementDate overriddenDate         `bson:"retirement_date"`
	Agencies       overriddenAgencies     `bson:"agencies"`
	Affiliations   overriddenAffiliations `bson:"affiliations"`
	Channels       overriddenChannels     `bson:"channels"`
}

type overriddenDate struct {
	Flag     bool       `bson:"flag"`
	OldValue *time.Time `bson:"old_value"`
	Value    *time.Time `bson:"value"`
}

type overriddenAgencies struct {
	Flag     bool     `bson:"flag"`
	OldValue []agency `bson:"old_value"`
	Value    []agency `bson:"value"`
}

type overriddenAffiliations struct {
	Flag     bool     `bson:"flag"`
	OldValue []string `bson:"old_value"`
	Value    []string `bson:"value"`
}

type overriddenChannels struct {
	Flag     bool      `bson:"flag"`
	OldValue []channel `bson:"old_value"`
	Value    []channel `bson:"value"`
}

func (o *overriddenField) toEntity() entity.OverriddenField {
	agencies := make([]entity.Agency, len(o.Agencies.Value))
	for i, a := range o.Agencies.Value {
		agencies[i] = entity.Agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	oldAgencies := make([]entity.Agency, len(o.Agencies.OldValue))
	for i, a := range o.Agencies.OldValue {
		oldAgencies[i] = entity.Agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]entity.Channel, len(o.Channels.Value))
	for i, c := range o.Channels.Value {
		channels[i] = entity.Channel{
			Type: c.Type,
			URL:  c.URL,
		}
	}

	oldChannels := make([]entity.Channel, len(o.Channels.OldValue))
	for i, c := range o.Channels.OldValue {
		oldChannels[i] = entity.Channel{
			Type: c.Type,
			URL:  c.URL,
		}
	}

	return entity.OverriddenField{
		DebutDate: entity.OverriddenDate{
			Flag:     o.DebutDate.Flag,
			OldValue: o.DebutDate.OldValue,
			Value:    o.DebutDate.Value,
		},
		RetirementDate: entity.OverriddenDate{
			Flag:     o.RetirementDate.Flag,
			OldValue: o.RetirementDate.OldValue,
			Value:    o.RetirementDate.Value,
		},
		Agencies: entity.OverriddenAgencies{
			Flag:     o.Agencies.Flag,
			OldValue: oldAgencies,
			Value:    agencies,
		},
		Affiliations: entity.OverriddenAffiliations{
			Flag:     o.Affiliations.Flag,
			OldValue: o.Affiliations.OldValue,
			Value:    o.Affiliations.Value,
		},
		Channels: entity.OverriddenChannels{
			Flag:     o.Channels.Flag,
			OldValue: oldChannels,
			Value:    channels,
		},
	}
}

func (m *Mongo) overiddenFieldFromEntity(o entity.OverriddenField) overriddenField {
	agencies := make([]agency, len(o.Agencies.Value))
	for i, a := range o.Agencies.Value {
		agencies[i] = agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	oldAgencies := make([]agency, len(o.Agencies.OldValue))
	for i, a := range o.Agencies.OldValue {
		oldAgencies[i] = agency{
			ID:    a.ID,
			Name:  a.Name,
			Image: a.Image,
		}
	}

	channels := make([]channel, len(o.Channels.Value))
	for i, c := range o.Channels.Value {
		channels[i] = channel{
			Type: c.Type,
			URL:  c.URL,
		}
	}

	oldChannels := make([]channel, len(o.Channels.OldValue))
	for i, c := range o.Channels.OldValue {
		oldChannels[i] = channel{
			Type: c.Type,
			URL:  c.URL,
		}
	}

	return overriddenField{
		DebutDate: overriddenDate{
			Flag:     o.DebutDate.Flag,
			OldValue: o.DebutDate.OldValue,
			Value:    o.DebutDate.Value,
		},
		RetirementDate: overriddenDate{
			Flag:     o.RetirementDate.Flag,
			OldValue: o.RetirementDate.OldValue,
			Value:    o.RetirementDate.Value,
		},
		Agencies: overriddenAgencies{
			Flag:     o.Agencies.Flag,
			OldValue: oldAgencies,
			Value:    agencies,
		},
		Affiliations: overriddenAffiliations{
			Flag:     o.Affiliations.Flag,
			OldValue: o.Affiliations.OldValue,
			Value:    o.Affiliations.Value,
		},
		Channels: overriddenChannels{
			Flag:     o.Channels.Flag,
			OldValue: oldChannels,
			Value:    channels,
		},
	}
}
