package mongo

import (
	"time"

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
