package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/language/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type language struct {
	ID        int64     `bson:"id"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// MarshalBSON to override marshal function.
func (l *language) MarshalBSON() ([]byte, error) {
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}

	l.UpdatedAt = time.Now()

	type l2 language
	return bson.Marshal((*l2)(l))
}

func (l *language) toEntity() *entity.Language {
	return &entity.Language{
		ID:   l.ID,
		Name: l.Name,
	}
}

func (m *Mongo) languageFromEntity(l entity.Language) *language {
	return &language{
		ID:   l.ID,
		Name: l.Name,
	}
}

func (m *Mongo) convertSort(sort string) bson.D {
	if sort == "" {
		sort = "name"
	}

	if sort[0] == '-' {
		return bson.D{{Key: sort[1:], Value: -1}}
	}

	return bson.D{{Key: sort, Value: 1}}
}
