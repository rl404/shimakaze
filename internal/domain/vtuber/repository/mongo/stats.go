package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCount to get count.
func (m *Mongo) GetCount(ctx context.Context) (int, int, error) {
	cnt, err := m.db.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return int(cnt), http.StatusOK, nil
}

// GetAverageActiveTime to get average active time.
func (m *Mongo) GetAverageActiveTime(ctx context.Context) (float64, int, error) {
	replaceFieldStage := bson.D{{Key: "$addFields", Value: bson.M{"new_retirement_date": bson.M{"$ifNull": []interface{}{"$retirement_date", primitive.NewDateTimeFromTime(time.Now())}}}}}
	matchStage := bson.D{{Key: "$match", Value: bson.M{"debut_date": bson.M{"$ne": nil}}}}
	avgStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": bson.M{"$dateDiff": bson.M{"startDate": "$debut_date", "endDate": "$new_retirement_date", "unit": "day"}}}}}}

	avgCursor, err := m.db.Aggregate(ctx, m.getPipeline(replaceFieldStage, matchStage, avgStage))
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return avg[0]["avg"], http.StatusOK, nil
}

// GetStatusCount to get status count.
func (m *Mongo) GetStatusCount(ctx context.Context) (*entity.StatusCount, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{
		"active":  bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$retirement_date", nil}}, 1, 0}},
		"retired": bson.M{"$cond": []interface{}{bson.M{"$ne": []interface{}{"$retirement_date", nil}}, 1, 0}}}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":     nil,
		"active":  bson.M{"$sum": "$active"},
		"retired": bson.M{"$sum": "$retired"},
	}}}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, groupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []map[string]int
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return &entity.StatusCount{
		Active:  cnt[0]["active"],
		Retired: cnt[0]["retired"],
	}, http.StatusOK, nil
}

type statusCountMonthly struct {
	Month int `bson:"month"`
	Year  int `bson:"year"`
	Count int `bson:"count"`
}

// GetDebutRetireCountMonthly to get debut & retire count monthly.
func (m *Mongo) GetDebutRetireCountMonthly(ctx context.Context) ([]entity.DebutRetireCount, int, error) {
	debutFilterStage := bson.D{{Key: "$match", Value: bson.M{"debut_date": bson.M{"$ne": nil}}}}
	debutProjectStage := bson.D{{Key: "$project", Value: bson.M{"month": bson.M{"$month": "$debut_date"}, "year": bson.M{"$year": "$debut_date"}}}}
	debutGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"month": "$month", "year": "$year"}, "count": bson.M{"$sum": 1}}}}
	debutProjectStage2 := bson.D{{Key: "$project", Value: bson.M{"month": "$_id.month", "year": "$_id.year", "count": "$count", "_id": 0}}}

	retiredFilterStage := bson.D{{Key: "$match", Value: bson.M{"retirement_date": bson.M{"$ne": nil}}}}
	retiredProjectStage := bson.D{{Key: "$project", Value: bson.M{"month": bson.M{"$month": "$retirement_date"}, "year": bson.M{"$year": "$retirement_date"}}}}
	retiredGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"month": "$month", "year": "$year"}, "count": bson.M{"$sum": 1}}}}
	retiredProjectStage2 := bson.D{{Key: "$project", Value: bson.M{"month": "$_id.month", "year": "$_id.year", "count": "$count", "_id": 0}}}

	debutCursor, err := m.db.Aggregate(ctx, m.getPipeline(debutFilterStage, debutProjectStage, debutGroupStage, debutProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	retiredCursor, err := m.db.Aggregate(ctx, m.getPipeline(retiredFilterStage, retiredProjectStage, retiredGroupStage, retiredProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var debutCount []statusCountMonthly
	if err := debutCursor.All(ctx, &debutCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var retiredCount []statusCountMonthly
	if err := retiredCursor.All(ctx, &retiredCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return m.mergeDebutRetiredMonthly(debutCount, retiredCount), http.StatusOK, nil
}

type statusCountYearly struct {
	Year  int `bson:"year"`
	Count int `bson:"count"`
}

// GetDebutRetireCountYearly to get debut & retire count yearly.
func (m *Mongo) GetDebutRetireCountYearly(ctx context.Context) ([]entity.DebutRetireCount, int, error) {
	debutFilterStage := bson.D{{Key: "$match", Value: bson.M{"debut_date": bson.M{"$ne": nil}}}}
	debutProjectStage := bson.D{{Key: "$project", Value: bson.M{"year": bson.M{"$year": "$debut_date"}}}}
	debutGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"year": "$year"}, "count": bson.M{"$sum": 1}}}}
	debutProjectStage2 := bson.D{{Key: "$project", Value: bson.M{"year": "$_id.year", "count": "$count", "_id": 0}}}

	retiredFilterStage := bson.D{{Key: "$match", Value: bson.M{"retirement_date": bson.M{"$ne": nil}}}}
	retiredProjectStage := bson.D{{Key: "$project", Value: bson.M{"year": bson.M{"$year": "$retirement_date"}}}}
	retiredGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"year": "$year"}, "count": bson.M{"$sum": 1}}}}
	retiredProjectStage2 := bson.D{{Key: "$project", Value: bson.M{"year": "$_id.year", "count": "$count", "_id": 0}}}

	debutCursor, err := m.db.Aggregate(ctx, m.getPipeline(debutFilterStage, debutProjectStage, debutGroupStage, debutProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	retiredCursor, err := m.db.Aggregate(ctx, m.getPipeline(retiredFilterStage, retiredProjectStage, retiredGroupStage, retiredProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var debutCount []statusCountYearly
	if err := debutCursor.All(ctx, &debutCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var retiredCount []statusCountYearly
	if err := retiredCursor.All(ctx, &retiredCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return m.mergeDebutRetiredYearly(debutCount, retiredCount), http.StatusOK, nil
}
