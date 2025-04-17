package mongo

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/non_vtuber/entity"
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
func (m *Mongo) Create(ctx context.Context, id int64, name string) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	if _, err := m.db.InsertOne(ctx, &nonVtuber{ID: id, Name: name}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusCreated, nil
}

// GetAllIDs to get all ids.
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	c, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var nonVtuber nonVtuber
		if err := c.Decode(&nonVtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		ids = append(ids, nonVtuber.ID)
	}

	return ids, http.StatusOK, nil
}

// GetAll to get list.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.NonVtuber, int, int, error) {
	query := bson.M{}
	opt := options.Find().SetSort(bson.D{{Key: "name", Value: 1}}).SetSkip(int64((data.Page - 1) * data.Limit)).SetLimit(int64(data.Limit))

	if data.Name != "" {
		query = bson.M{"name": bson.M{"$regex": data.Name, "$options": "i"}}
	}

	if data.Limit < 0 {
		opt.SetLimit(0)
	}

	c, err := m.db.Find(ctx, query, opt)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	var nonVtubers []entity.NonVtuber
	for c.Next(ctx) {
		var nonVtuber nonVtuber
		if err := c.Decode(&nonVtuber); err != nil {
			return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		nonVtubers = append(nonVtubers, nonVtuber.toEntity())
	}

	total, err := m.db.CountDocuments(ctx, bson.M{}, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return nonVtubers, int(total), http.StatusOK, nil
}

// DeleteByID to delete by id.
func (m *Mongo) DeleteByID(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}
