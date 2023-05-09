package mongo

import (
	"context"
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for non-vtuber mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new non-vtuber mongodb.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("non_vtuber"),
	}
}

// Create to create non-vtuber.
func (m *Mongo) Create(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	if _, err := m.db.InsertOne(ctx, &nonVtuber{ID: id}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return http.StatusCreated, nil
}

// GetAllIDs
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	c, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var nonVtuber nonVtuber
		if err := c.Decode(&nonVtuber); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
		ids = append(ids, nonVtuber.ID)
	}

	return ids, http.StatusOK, nil
}
