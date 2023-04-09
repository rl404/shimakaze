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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []designerCount
	if err := countCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var avg []map[string]float64
	if err := avgCursor.All(ctx, &avg); err != nil {
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []videoCountByDate
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var cnt []videoCount
	if err := cntCursor.All(ctx, &cnt); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var dur []videoDuration
	if err := durCursor.All(ctx, &dur); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
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
