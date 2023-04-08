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
	replaceFieldStage := bson.D{{Key: "$addFields", Value: bson.M{"new_retirement_date": bson.M{"$ifNull": bson.A{"$retirement_date", primitive.NewDateTimeFromTime(time.Now())}}}}}
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
		"active":  bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$retirement_date", nil}}, 1, 0}},
		"retired": bson.M{"$cond": bson.A{bson.M{"$ne": bson.A{"$retirement_date", nil}}, 1, 0}}}}}

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

type modelCount struct {
	None      int `bson:"none"`
	Has2DOnly int `bson:"has_2d_only"`
	Has3DOnly int `bson:"has_3d_only"`
	Both      int `bson:"both"`
}

// GetModelCount to get 2d & 3d model count.
func (m *Mongo) GetModelCount(ctx context.Context) (*entity.ModelCount, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{
		"none":        bson.M{"$cond": bson.A{bson.M{"$and": []bson.M{{"$eq": bson.A{"$has_2d", false}}, {"$eq": bson.A{"$has_3d", false}}}}, 1, 0}},
		"has_2d_only": bson.M{"$cond": bson.A{bson.M{"$and": []bson.M{{"$eq": bson.A{"$has_2d", true}}, {"$eq": bson.A{"$has_3d", false}}}}, 1, 0}},
		"has_3d_only": bson.M{"$cond": bson.A{bson.M{"$and": []bson.M{{"$eq": bson.A{"$has_2d", false}}, {"$eq": bson.A{"$has_3d", true}}}}, 1, 0}},
		"both":        bson.M{"$cond": bson.A{bson.M{"$and": []bson.M{{"$eq": bson.A{"$has_2d", true}}, {"$eq": bson.A{"$has_3d", true}}}}, 1, 0}},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":         nil,
		"none":        bson.M{"$sum": "$none"},
		"has_2d_only": bson.M{"$sum": "$has_2d_only"},
		"has_3d_only": bson.M{"$sum": "$has_3d_only"},
		"both":        bson.M{"$sum": "$both"},
	}}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, groupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []modelCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return &entity.ModelCount{
		None:      cnt[0].None,
		Has2DOnly: cnt[0].Has2DOnly,
		Has3DOnly: cnt[0].Has3DOnly,
		Both:      cnt[0].Both,
	}, http.StatusOK, nil
}

type inAgencyCount struct {
	InAgency    int `bson:"in_agency"`
	NotInAgency int `bson:"not_in_agency"`
}

// GetInAgencyCount to get in agency count.
func (m *Mongo) GetInAgencyCount(ctx context.Context) (*entity.InAgencyCount, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{
		"in_agency": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"array", bson.M{"$type": "$agencies"}}},
			bson.M{"$cond": bson.A{bson.M{"$gt": bson.A{bson.M{"$size": "$agencies"}, 0}}, 1, 0}},
			0}},
		"not_in_agency": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"array", bson.M{"$type": "$agencies"}}},
			bson.M{"$cond": bson.A{bson.M{"$gt": bson.A{bson.M{"$size": "$agencies"}, 0}}, 0, 1}},
			1}},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":           nil,
		"in_agency":     bson.M{"$sum": "$in_agency"},
		"not_in_agency": bson.M{"$sum": "$not_in_agency"},
	}}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, groupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []inAgencyCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return &entity.InAgencyCount{
		InAgency:    cnt[0].InAgency,
		NotInAgency: cnt[0].NotInAgency,
	}, http.StatusOK, nil
}

type subscriberCount struct {
	Min   int `bson:"min"`
	Max   int `bson:"max"`
	Count int `bson:"count"`
}

// GetSubscriberCount to get subscriber count.
func (m *Mongo) GetSubscriberCount(ctx context.Context, interval, max int) ([]entity.SubscriberCount, int, error) {
	boundaries := []int{}
	for i := 0; i <= max; i += interval {
		boundaries = append(boundaries, i)
	}

	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{"subscriber": bson.M{"$max": "$channels.subscriber"}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.M{"subscriber": bson.M{"$ifNull": bson.A{"$subscriber", 0}}}}}

	bucketStage := bson.D{{Key: "$bucket", Value: bson.M{
		"groupBy":    "$subscriber",
		"boundaries": boundaries,
		"default":    max,
		"output":     bson.M{"count": bson.M{"$sum": 1}},
	}}}

	projectStage2 := bson.D{{Key: "$project", Value: bson.M{
		"min":   "$_id",
		"max":   bson.M{"$add": bson.A{"$_id", interval}},
		"count": "$count",
	}}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, projectStage, bucketStage, projectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []subscriberCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]entity.SubscriberCount, len(cnt))
	for i, c := range cnt {
		max := c.Max
		if i == len(cnt)-1 {
			max = 0
		}

		res[i] = entity.SubscriberCount{
			Min:   c.Min,
			Max:   max,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}
