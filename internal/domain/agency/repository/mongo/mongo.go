package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for agency mongodb.
type Mongo struct {
	db     *mongo.Collection
	oldAge time.Duration
}

// New to create new agency mongodb.
func New(db *mongo.Database, oldAge int) *Mongo {
	return &Mongo{
		db:     db.Collection("agency"),
		oldAge: time.Duration(oldAge) * 24 * time.Hour,
	}
}

// GetAllIDs to get all ids.
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	c, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var agency agency
		if err := c.Decode(&agency); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
		ids = append(ids, agency.ID)
	}

	return ids, http.StatusOK, nil
}

// GetAll to get all.
func (m *Mongo) GetAll(ctx context.Context) ([]entity.Agency, int, error) {
	var agencies []entity.Agency
	c, err := m.db.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var agency agency
		if err := c.Decode(&agency); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
		agencies = append(agencies, agency.toEntity())
	}

	return agencies, http.StatusOK, nil
}

// IsOld to check if old data.
func (m *Mongo) IsOld(ctx context.Context, id int64) (bool, int, error) {
	filter := bson.M{
		"id":         id,
		"updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.oldAge))},
	}

	if err := m.db.FindOne(ctx, filter).Decode(&agency{}); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return true, http.StatusNotFound, nil
		}
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return false, http.StatusOK, nil
}

// UpdateByID to update by id.
func (m *Mongo) UpdateByID(ctx context.Context, id int64, data entity.Agency) (int, error) {
	var agency agency
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&agency); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := m.db.InsertOne(ctx, m.agencyFromEntity(data)); err != nil {
				return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	mm := m.agencyFromEntity(data)
	mm.CreatedAt = agency.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": mm}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}
