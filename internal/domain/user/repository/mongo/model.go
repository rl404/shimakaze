package mongo

import (
	"time"

	"github.com/rl404/shimakaze/internal/domain/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type user struct {
	ID        int64     `bson:"id"`
	Username  string    `bson:"username"`
	IsAdmin   bool      `bson:"is_admin"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// MarshalBSON to override marshal function.
func (n *user) MarshalBSON() ([]byte, error) {
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	n.UpdatedAt = time.Now()

	type n2 user
	return bson.Marshal((*n2)(n))
}

func (n *user) toEntity() *entity.User {
	return &entity.User{
		ID:       n.ID,
		Username: n.Username,
		IsAdmin:  n.IsAdmin,
	}
}

func (m *Mongo) fromEntity(data entity.User) *user {
	return &user{
		ID:       data.ID,
		Username: data.Username,
		IsAdmin:  data.IsAdmin,
	}
}
