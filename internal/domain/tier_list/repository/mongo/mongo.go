package mongo

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/tier_list/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo contains function for tier list mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new tier list mongodb.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("tier_list"),
	}
}

// Get to get list.
func (m *Mongo) Get(ctx context.Context, data entity.GetRequest) ([]entity.TierList, int, int, error) {
	matchStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: m.convertSort(data.Sort)}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.Query != "" {
		matchStage = m.addMatch(matchStage, "$or", []bson.M{
			{"title": bson.M{"$regex": data.Query, "$options": "i"}},
			{"description": bson.M{"$regex": data.Query, "$options": "i"}},
		})
	}

	if data.UserID != 0 {
		matchStage = m.addMatch(matchStage, "created_by.id", data.UserID)
	}

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, sortStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var tierlists []tierList
	if err := cursor.All(ctx, &tierlists); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]entity.TierList, len(tierlists))
	for i, t := range tierlists {
		res[i] = *t.toEntity()
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, countStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var total []map[string]int64
	if err := cntCursor.All(ctx, &total); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if len(total) == 0 {
		return res, 0, http.StatusOK, nil
	}

	return res, int(total[0]["count"]), http.StatusOK, nil
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id string) (*entity.TierList, int, error) {
	var tierList tierList
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&tierList); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrTierNotFound)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return tierList.toEntity(), http.StatusOK, nil
}

// UpsertByID to upsert by id.
func (m *Mongo) UpsertByID(ctx context.Context, data entity.TierList) (*entity.TierList, int, error) {
	if data.ID == "" {
		tierList := m.fromEntity(data)
		if _, err := m.db.InsertOne(ctx, tierList); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		return tierList.toEntity(), http.StatusCreated, nil
	}

	var tierList tierList
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&tierList); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrTierNotFound)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if tierList.CreatedBy.ID != data.User.ID {
		return nil, http.StatusForbidden, stack.Wrap(ctx, errors.ErrUpdateNotAllowed)
	}

	tl := m.fromEntity(data)
	tl.CreatedAt = tierList.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": tl}); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return tl.toEntity(), http.StatusOK, nil
}

// DeleteByID to delete by id.
func (m *Mongo) DeleteByID(ctx context.Context, id string, userID int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}
