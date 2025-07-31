package mongo

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/language/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Mongo contains functions for language mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new language mongodb.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("language"),
	}
}

// GetAllIDs to get all by ids.
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	c, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var language language
		if err := c.Decode(&language); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		ids = append(ids, language.ID)
	}

	return ids, http.StatusOK, nil
}

// UpdateByID to update by id.
func (m *Mongo) UpdateByID(ctx context.Context, id int64, data entity.Language) (int, error) {
	var language language
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&language); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := m.db.InsertOne(ctx, m.languageFromEntity(data)); err != nil {
				return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	mm := m.languageFromEntity(data)
	mm.CreatedAt = language.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": mm}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// GetAll to get all.
func (m *Mongo) GetAll(ctx context.Context) ([]entity.Language, int, int, error) {
	c, err := m.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	var languages []entity.Language
	for c.Next(ctx) {
		var language language
		if err := c.Decode(&language); err != nil {
			return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		languages = append(languages, *language.toEntity())
	}

	total, err := m.db.CountDocuments(ctx, bson.M{}, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return languages, int(total), http.StatusOK, nil
}
