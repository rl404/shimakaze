package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/non_vtuber/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type nonVtuber struct {
	ID        int64     `bson:"id"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}

// MarshalBSON to override marshal function.
func (n *nonVtuber) MarshalBSON() ([]byte, error) {
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	type n2 nonVtuber
	return bson.Marshal((*n2)(n))
}

func (n *nonVtuber) toEntity() entity.NonVtuber {
	return entity.NonVtuber{
		ID:   n.ID,
		Name: n.Name,
	}
}
