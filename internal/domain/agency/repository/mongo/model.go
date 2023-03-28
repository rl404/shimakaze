package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type agency struct {
	ID        int64     `bson:"id"`
	Name      string    `bson:"name"`
	Image     string    `bson:"image"`
	Member    int       `bson:"member"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// MarshalBSON to override marshal function.
func (a *agency) MarshalBSON() ([]byte, error) {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}

	a.UpdatedAt = time.Now()

	type a2 agency
	return bson.Marshal((*a2)(a))
}

func (a *agency) toEntity() *entity.Agency {
	return &entity.Agency{
		ID:        a.ID,
		Name:      a.Name,
		Image:     a.Image,
		Member:    a.Member,
		UpdatedAt: a.UpdatedAt,
	}
}

func (m *Mongo) agencyFromEntity(a entity.Agency) *agency {
	return &agency{
		ID:        a.ID,
		Name:      a.Name,
		Image:     a.Image,
		Member:    a.Member,
		UpdatedAt: a.UpdatedAt,
	}
}
