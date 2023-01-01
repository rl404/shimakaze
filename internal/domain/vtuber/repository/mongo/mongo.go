package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for vtuber mongodb.
type Mongo struct {
	db     *mongo.Collection
	oldAge time.Duration
}

// New to create new vtuber mongodb.
func New(db *mongo.Database, oldAge int) *Mongo {
	return &Mongo{
		db:     db.Collection("vtuber"),
		oldAge: time.Duration(oldAge) * 24 * time.Hour,
	}
}

// GetAllIDs to get all ids.
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var ids []int64
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		ids = append(ids, vtuber.ID)
	}

	return ids, http.StatusOK, nil
}

// UpdateByID to update by id.
func (m *Mongo) UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error) {
	var vtuber vtuber
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&vtuber); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := m.db.InsertOne(ctx, m.vtuberFromEntity(data)); err != nil {
				return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	mm := m.vtuberFromEntity(data)
	mm.CreatedAt = vtuber.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": mm}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}

// IsOld to check if old data.
func (m *Mongo) IsOld(ctx context.Context, id int64) (bool, int, error) {
	filter := bson.M{
		"id":         id,
		"updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.oldAge))},
	}

	if err := m.db.FindOne(ctx, filter).Decode(&vtuber{}); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return true, http.StatusNotFound, nil
		}
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return false, http.StatusOK, nil
}
