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

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Vtuber, int, error) {
	var vtuber vtuber
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&vtuber); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrVtuberNotFound, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return vtuber.toEntity(), http.StatusOK, nil
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

// GetOldIDs to get old ids.
func (m *Mongo) GetOldIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{
		"updated_at": bson.M{"$lte": primitive.NewDateTimeFromTime(time.Now().Add(-m.oldAge))},
	}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer cursor.Close(ctx)

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

// GetAllImages to get all images.
func (m *Mongo) GetAllImages(ctx context.Context, _ bool, _ int) ([]entity.Vtuber, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{"image": bson.M{"$ne": ""}}, options.Find().SetProjection(bson.M{"id": 1, "name": 1, "image": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		res = append(res, entity.Vtuber{
			ID:    vtuber.ID,
			Name:  vtuber.Name,
			Image: vtuber.Image,
		})
	}

	return res, http.StatusOK, nil
}

// GetAllForFamilyTree to get all data for tree.
func (m *Mongo) GetAllForFamilyTree(ctx context.Context) ([]entity.Vtuber, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"id":                    1,
		"name":                  1,
		"image":                 1,
		"retirement_date":       1,
		"character_designers":   1,
		"character_2d_modelers": 1,
		"character_3d_modelers": 1,
	}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		res = append(res, entity.Vtuber{
			ID:                  vtuber.ID,
			Name:                vtuber.Name,
			Image:               vtuber.Image,
			RetirementDate:      vtuber.RetirementDate,
			CharacterDesigners:  vtuber.CharacterDesigners,
			Character2DModelers: vtuber.Character2DModelers,
			Character3DModelers: vtuber.Character3DModelers,
		})
	}

	return res, http.StatusOK, nil
}

// GetAllForAgencyTree to get all data for agency tree.
func (m *Mongo) GetAllForAgencyTree(ctx context.Context) ([]entity.Vtuber, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"id":              1,
		"name":            1,
		"image":           1,
		"retirement_date": 1,
		"agencies":        1,
	}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		agencies := make([]entity.Agency, len(vtuber.Agencies))
		for i, a := range vtuber.Agencies {
			agencies[i] = entity.Agency{
				ID:   a.ID,
				Name: a.Name,
			}
		}

		res = append(res, entity.Vtuber{
			ID:             vtuber.ID,
			Name:           vtuber.Name,
			Image:          vtuber.Image,
			RetirementDate: vtuber.RetirementDate,
			Agencies:       agencies,
		})
	}

	return res, http.StatusOK, nil
}

// GetAll to get all data.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Vtuber, int, int, error) {
	filter := bson.M{}
	opt := options.Find().SetSkip(int64((data.Page - 1) * data.Limit)).SetLimit(int64(data.Limit))

	if data.Mode == entity.SearchModeStats {
		opt.SetProjection(bson.M{
			"image":             0,
			"original_names":    0,
			"nicknames":         0,
			"caption":           0,
			"affiliations":      0,
			"official_websites": 0,
		})
	}

	if data.Limit < 0 {
		opt.SetLimit(0)
	}

	cursor, err := m.db.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
		res = append(res, *vtuber.toEntity())
	}

	total, err := m.db.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return res, int(total), http.StatusOK, nil
}
