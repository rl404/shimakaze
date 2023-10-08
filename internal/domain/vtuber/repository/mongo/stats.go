package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	_errors "github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCount to get count.
func (m *Mongo) GetCount(ctx context.Context) (int, int, error) {
	cnt, err := m.db.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []map[string]int
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	retiredCursor, err := m.db.Aggregate(ctx, m.getPipeline(retiredFilterStage, retiredProjectStage, retiredGroupStage, retiredProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var debutCount []statusCountMonthly
	if err := debutCursor.All(ctx, &debutCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var retiredCount []statusCountMonthly
	if err := retiredCursor.All(ctx, &retiredCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	retiredCursor, err := m.db.Aggregate(ctx, m.getPipeline(retiredFilterStage, retiredProjectStage, retiredGroupStage, retiredProjectStage2))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var debutCount []statusCountYearly
	if err := debutCursor.All(ctx, &debutCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var retiredCount []statusCountYearly
	if err := retiredCursor.All(ctx, &retiredCount); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []modelCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []inAgencyCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []subscriberCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
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

type designerCount struct {
	Name  string `bson:"name"`
	Count int    `bson:"count"`
}

func (m *Mongo) getdesignerCount(ctx context.Context, top int, field string) ([]entity.DesignerCount, int, error) {
	unwindStage := bson.D{{Key: "$unwind", Value: "$" + field}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"name": "$" + field}, "count": bson.M{"$sum": 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"name": "$_id.name", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: top}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(unwindStage, groupStage, projectStage, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []designerCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.DesignerCount, len(cnt))
	for i, c := range cnt {
		res[i] = entity.DesignerCount{
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetDesignerCount to get character designer count.
func (m *Mongo) GetDesignerCount(ctx context.Context, top int) ([]entity.DesignerCount, int, error) {
	return m.getdesignerCount(ctx, top, "character_designers")
}

// Get2DModelerCount to get character 2d modeler count.
func (m *Mongo) Get2DModelerCount(ctx context.Context, top int) ([]entity.DesignerCount, int, error) {
	return m.getdesignerCount(ctx, top, "character_2d_modelers")
}

// Get2DModelerCount to get character 3d modeler count.
func (m *Mongo) Get3DModelerCount(ctx context.Context, top int) ([]entity.DesignerCount, int, error) {
	return m.getdesignerCount(ctx, top, "character_3d_modelers")
}

// GetAverageVideoCount to get average video count.
func (m *Mongo) GetAverageVideoCount(ctx context.Context) (float64, int, error) {
	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{
		"channels": bson.M{"$map": bson.M{
			"input": "$channels",
			"as":    "channel",
			"in": bson.M{
				"$mergeObjects": bson.A{"$$channel", bson.M{
					"video_count": bson.M{"$size": "$$channel.videos"},
				}},
			},
		}}}}}
	newFieldStage2 := bson.D{{Key: "$addFields", Value: bson.M{"video_count": bson.M{"$sum": "$channels.video_count"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"name": 1, "video_count": 1}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$video_count"}}}}

	avgCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, newFieldStage2, projectStage, groupStage))
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	return avg[0]["avg"], http.StatusOK, nil
}

// GetAverageVideoDuration to get average video duration.
func (m *Mongo) GetAverageVideoDuration(ctx context.Context) (float64, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{"videos": "$channels.videos"}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{"path": "$videos"}}}
	matchStage := bson.D{{Key: "$match", Value: bson.M{"videos.end_date": bson.M{"$ne": nil}}}}
	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{"duration": bson.M{
		"$dateDiff": bson.M{
			"startDate": "$videos.start_date",
			"endDate":   "$videos.end_date",
			"unit":      "second",
		},
	}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$duration"}}}}

	avgCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, unwindStage, unwindStage, matchStage, newFieldStage, groupStage))
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	return avg[0]["avg"], http.StatusOK, nil
}

type videoCountByDate struct {
	Day   int `bson:"day"`
	Hour  int `bson:"hour"`
	Count int `bson:"count"`
}

// GetVideoCountByDate to get video count by date.
func (m *Mongo) GetVideoCountByDate(ctx context.Context, hourly, daily bool) ([]entity.VideoCountByDate, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{"videos": "$channels.videos"}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{"path": "$videos"}}}
	matchStage := bson.D{{Key: "$match", Value: bson.M{"videos.start_date": bson.M{"$ne": nil}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{}, "count": bson.M{"$sum": 1}}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.M{"day": "$_id.day", "hour": "$_id.hour", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"day": 1, "hour": 1}}}

	if hourly {
		groupStage[0].Value.(bson.M)["_id"].(bson.M)["hour"] = bson.M{"$hour": "$videos.start_date"}
	}

	if daily {
		groupStage[0].Value.(bson.M)["_id"].(bson.M)["day"] = bson.M{"$dayOfWeek": "$videos.start_date"}
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, unwindStage, unwindStage, matchStage, groupStage, projectStage2, sortStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []videoCountByDate
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.VideoCountByDate, len(cnt))
	for i, c := range cnt {
		res[i] = entity.VideoCountByDate{
			Day:   c.Day,
			Hour:  c.Hour,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type videoCount struct {
	ID    int64  `bson:"id"`
	Name  string `bson:"name"`
	Count int    `bson:"count"`
}

// GetVideoCount to get video count.
func (m *Mongo) GetVideoCount(ctx context.Context, top int) ([]entity.VideoCount, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{"id": 1, "name": 1, "videos": "$channels.videos"}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{"path": "$videos"}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"id": "$id", "name": "$name"}, "count": bson.M{"$sum": 1}}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.M{"id": "$_id.id", "name": "$_id.name", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: top}}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, unwindStage, unwindStage, groupStage, projectStage2, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []videoCount
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.VideoCount, len(cnt))
	for i, c := range cnt {
		res[i] = entity.VideoCount{
			ID:    c.ID,
			Name:  c.Name,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type videoDuration struct {
	ID       int64   `bson:"id"`
	Name     string  `bson:"name"`
	Duration float64 `bson:"duration"`
}

// GetVideoDuration to get video duration.
func (m *Mongo) GetVideoDuration(ctx context.Context, top int) ([]entity.VideoDuration, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{"id": 1, "name": 1, "videos": "$channels.videos"}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{"path": "$videos"}}}
	matchStage := bson.D{{Key: "$match", Value: bson.M{"videos.start_date": bson.M{"$ne": nil}, "videos.end_date": bson.M{"$ne": nil}}}}
	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{"duration": bson.M{
		"$dateDiff": bson.M{
			"startDate": "$videos.start_date",
			"endDate":   "$videos.end_date",
			"unit":      "second",
		},
	}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"id": "$id", "name": "$name"}, "duration": bson.M{"$avg": "$duration"}}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.M{"id": "$_id.id", "name": "$_id.name", "duration": "$duration"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"duration": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: top}}

	durCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, unwindStage, unwindStage, matchStage, newFieldStage, groupStage, projectStage2, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var dur []videoDuration
	if err := durCursor.All(ctx, &dur); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.VideoDuration, len(dur))
	for i, c := range dur {
		res[i] = entity.VideoDuration{
			ID:       c.ID,
			Name:     c.Name,
			Duration: c.Duration,
		}
	}

	return res, http.StatusOK, nil
}

type birthdayCount struct {
	Month int `bson:"month"`
	Day   int `bson:"day"`
	Count int `bson:"count"`
}

// GetBirthdayCount to get birthday count.
func (m *Mongo) GetBirthdayCount(ctx context.Context) ([]entity.BirthdayCount, int, error) {
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{
		"month": bson.M{"$month": "$birthday"},
		"day":   bson.M{"$dayOfMonth": "$birthday"},
	}, "count": bson.M{"$sum": 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"month": "$_id.month", "day": "$_id.day", "count": "$count"}}}
	matchStage := bson.D{{Key: "$match", Value: bson.M{"month": bson.M{"$gt": 0}, "day": bson.M{"$gt": 0}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"month": 1, "day": 1}}}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(groupStage, projectStage, matchStage, sortStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []birthdayCount
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.BirthdayCount, len(cnt))
	for i, c := range cnt {
		res[i] = entity.BirthdayCount{
			Month: c.Month,
			Day:   c.Day,
			Count: c.Count,
		}
	}

	return res, http.StatusOK, nil
}

// GetAverageHeight to get average height.
func (m *Mongo) GetAverageHeight(ctx context.Context) (float64, int, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"height": bson.M{"$gte": 0}}}}
	avgStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$height"}}}}

	avgCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, avgStage))
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	return avg[0]["avg"], http.StatusOK, nil
}

// GetAverageWeight to get average weight.
func (m *Mongo) GetAverageWeight(ctx context.Context) (float64, int, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"weight": bson.M{"$gte": 0, "$lte": 1000}}}}
	avgStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$weight"}}}}

	avgCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, avgStage))
	if err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	return avg[0]["avg"], http.StatusOK, nil
}

type bloodTypeCount struct {
	BloodType string `bson:"blood_type"`
	Count     int    `bson:"count"`
}

// GetBloodTypeCount to get blood type count.
func (m *Mongo) GetBloodTypeCount(ctx context.Context, top int) ([]entity.BloodTypeCount, int, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"blood_type": bson.M{"$nin": bson.A{"", nil}}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"blood_type": "$blood_type"}, "count": bson.M{"$sum": 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"blood_type": "$_id.blood_type", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: top}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, groupStage, projectStage, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []bloodTypeCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	topBloodTypes := make([]string, top)
	res := make([]entity.BloodTypeCount, len(cnt)+1)
	for i, c := range cnt {
		res[i] = entity.BloodTypeCount{
			BloodType: c.BloodType,
			Count:     c.Count,
		}
		topBloodTypes[i] = c.BloodType
	}

	// Get other blood type.
	otherMatchStage := bson.D{{Key: "$match", Value: bson.M{"blood_type": bson.M{"$nin": topBloodTypes}}}}
	otherGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}

	otherCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, otherMatchStage, otherGroupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var otherCnt []bloodTypeCount
	if err := otherCursor.All(ctx, &otherCnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res[top] = entity.BloodTypeCount{
		BloodType: "other",
		Count:     otherCnt[0].Count,
	}

	return res, http.StatusOK, nil
}

type channelTypeCount struct {
	ChannelType entity.ChannelType `bson:"channel_type"`
	Count       int                `bson:"count"`
}

// GetChannelTypeCount to get channel type count.
func (m *Mongo) GetChannelTypeCount(ctx context.Context) ([]entity.ChannelTypeCount, int, error) {
	projectStage := bson.D{{Key: "$project", Value: bson.M{"channels.type": 1}}}
	unwindStage := bson.D{{Key: "$unwind", Value: "$channels"}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"type": "$channels.type"}, "count": bson.M{"$sum": 1}}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.M{"channel_type": "$_id.type", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(projectStage, unwindStage, groupStage, projectStage2, sortStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []channelTypeCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res := make([]entity.ChannelTypeCount, len(cnt))
	for i, c := range cnt {
		res[i] = entity.ChannelTypeCount{
			ChannelType: c.ChannelType,
			Count:       c.Count,
		}
	}

	return res, http.StatusOK, nil
}

type genderCount struct {
	Gender string `bson:"gender"`
	Count  int    `bson:"count"`
}

// GetGenderCount to get gender count.
func (m *Mongo) GetGenderCount(ctx context.Context) ([]entity.GenderCount, int, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"gender": bson.M{"$nin": bson.A{"", nil}}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"gender": "$gender"}, "count": bson.M{"$sum": 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"gender": "$_id.gender", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: 2}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, groupStage, projectStage, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []genderCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	topGenders := make([]string, 2)
	res := make([]entity.GenderCount, 3)
	for i, c := range cnt {
		res[i] = entity.GenderCount{
			Gender: c.Gender,
			Count:  c.Count,
		}
		topGenders[i] = c.Gender
	}

	// Get other gender.
	otherMatchStage := bson.D{{Key: "$match", Value: bson.M{"gender": bson.M{"$nin": topGenders}}}}
	otherGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}

	otherCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, otherMatchStage, otherGroupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var otherCnt []genderCount
	if err := otherCursor.All(ctx, &otherCnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res[2] = entity.GenderCount{
		Gender: "other",
		Count:  otherCnt[0].Count,
	}

	return res, http.StatusOK, nil
}

type zodiacCount struct {
	Zodiac string `bson:"zodiac"`
	Count  int    `bson:"count"`
}

// GetZodiacCount to get zodiac count.
func (m *Mongo) GetZodiacCount(ctx context.Context) ([]entity.ZodiacCount, int, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"zodiac_sign": bson.M{"$nin": bson.A{"", nil}}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": bson.M{"zodiac": "$zodiac_sign"}, "count": bson.M{"$sum": 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{"zodiac": "$_id.zodiac", "count": "$count"}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"count": -1}}}
	limitStage := bson.D{{Key: "$limit", Value: 12}}

	countCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, groupStage, projectStage, sortStage, limitStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var cnt []zodiacCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	topZodiacs := make([]string, 12)
	res := make([]entity.ZodiacCount, 13)
	for i, c := range cnt {
		res[i] = entity.ZodiacCount{
			Zodiac: c.Zodiac,
			Count:  c.Count,
		}
		topZodiacs[i] = c.Zodiac
	}

	// Get other zodiac.
	otherMatchStage := bson.D{{Key: "$match", Value: bson.M{"zodiac_sign": bson.M{"$nin": topZodiacs}}}}
	otherGroupStage := bson.D{{Key: "$group", Value: bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}}

	otherCursor, err := m.db.Aggregate(ctx, m.getPipeline(matchStage, otherMatchStage, otherGroupStage))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	var otherCnt []genderCount
	if err := otherCursor.All(ctx, &otherCnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, err, _errors.ErrInternalDB)
	}

	res[12] = entity.ZodiacCount{
		Zodiac: "other",
		Count:  otherCnt[0].Count,
	}

	return res, http.StatusOK, nil
}
