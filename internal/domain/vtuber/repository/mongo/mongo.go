package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Mongo contains functions for vtuber mongodb.
type Mongo struct {
	db            *mongo.Collection
	oldActiveAge  time.Duration
	oldRetiredAge time.Duration
}

// New to create new vtuber mongodb.
func New(db *mongo.Database, oldActiveAge, oldRetiredAge int) *Mongo {
	return &Mongo{
		db:            db.Collection("vtuber"),
		oldActiveAge:  time.Duration(oldActiveAge) * 24 * time.Hour,
		oldRetiredAge: time.Duration(oldRetiredAge) * 24 * time.Hour,
	}
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Vtuber, int, error) {
	var vtuber vtuber
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&vtuber); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrVtuberNotFound)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return vtuber.toEntity(), http.StatusOK, nil
}

// GetAllIDs to get all ids.
func (m *Mongo) GetAllIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var ids []int64
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
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
				return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	mm := m.vtuberFromEntity(data)
	mm.CreatedAt = vtuber.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": mm}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// DeleteByID to delete by id.
func (m *Mongo) DeleteByID(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// IsOld to check if old data.
func (m *Mongo) IsOld(ctx context.Context, id int64) (bool, int, error) {
	filter := bson.M{
		"id": id,
		"$or": bson.A{
			bson.M{"retirement_date": bson.M{"$eq": nil}, "updated_at": bson.M{"$gte": bson.NewDateTimeFromTime(time.Now().Add(-m.oldActiveAge))}},
			bson.M{"retirement_date": bson.M{"$ne": nil}, "updated_at": bson.M{"$gte": bson.NewDateTimeFromTime(time.Now().Add(-m.oldRetiredAge))}},
		},
	}

	if err := m.db.FindOne(ctx, filter).Decode(&vtuber{}); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return true, http.StatusNotFound, nil
		}
		return true, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return false, http.StatusOK, nil
}

// GetOldActiveIDs to get old active ids.
func (m *Mongo) GetOldActiveIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{
		"retirement_date": bson.M{"$eq": nil},
		"updated_at":      bson.M{"$lte": bson.NewDateTimeFromTime(time.Now().Add(-m.oldActiveAge))},
	}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer cursor.Close(ctx)

	var ids []int64
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}

		ids = append(ids, vtuber.ID)
	}

	return ids, http.StatusOK, nil
}

// GetOldRetiredIDs to get old retired ids.
func (m *Mongo) GetOldRetiredIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{
		"retirement_date": bson.M{"$ne": nil},
		"updated_at":      bson.M{"$lte": bson.NewDateTimeFromTime(time.Now().Add(-m.oldActiveAge))},
	}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer cursor.Close(ctx)

	var ids []int64
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}

		ids = append(ids, vtuber.ID)
	}

	return ids, http.StatusOK, nil
}

// GetAllImages to get all images.
func (m *Mongo) GetAllImages(ctx context.Context, _ bool, _ int) ([]entity.Vtuber, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{"image": bson.M{"$ne": ""}}, options.Find().SetProjection(bson.M{"id": 1, "name": 1, "image": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
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
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var res []entity.Vtuber
	for cursor.Next(ctx) {
		var vtuber vtuber
		if err := cursor.Decode(&vtuber); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
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
	newFieldStage := bson.D{}
	matchStage := bson.D{}
	projectStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: m.convertSort(data.Sort)}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.Mode == entity.SearchModeSimple {
		projectStage = bson.D{{Key: "$project", Value: bson.M{
			"id":                   1,
			"name":                 1,
			"image":                1,
			"debut_date":           1,
			"retirement_date":      1,
			"subscriber":           1,
			"monthly_subscriber":   1,
			"video_count":          1,
			"average_video_length": 1,
			"total_video_length":   1,
			"has_2d":               1,
			"has_3d":               1,
			"agencies":             1,
			"birthday":             1,
			"emoji":                1,
			"updated_at":           1,
			"is_debut_date_null": bson.M{"$cond": bson.A{
				bson.M{"$eq": bson.A{"$debut_date", nil}},
				1, 0,
			}},
		}}}
	}

	if data.Names != "" {
		matchStage = m.addMatch(matchStage, "$or", []bson.M{
			{"name": bson.M{"$regex": data.Names, "$options": "i"}},
			{"original_names": bson.M{"$regex": data.Names, "$options": "i"}},
			{"nicknames": bson.M{"$regex": data.Names, "$options": "i"}},
		})
	}

	if data.Name != "" {
		matchStage = m.addMatch(matchStage, "name", bson.M{"$regex": data.Name, "$options": "i"})
	}

	if data.OriginalName != "" {
		matchStage = m.addMatch(matchStage, "original_names", bson.M{"$regex": data.OriginalName, "$options": "i"})
	}

	if data.Nickname != "" {
		matchStage = m.addMatch(matchStage, "nicknames", bson.M{"$regex": data.Nickname, "$options": "i"})
	}

	if data.ExcludeActive {
		matchStage = m.addMatch(matchStage, "retirement_date", bson.M{"$ne": nil})
	}

	if data.ExcludeRetired {
		matchStage = m.addMatch(matchStage, "retirement_date", bson.M{"$eq": nil})
	}

	if data.ExcludeActive && data.ExcludeRetired {
		return nil, 0, http.StatusOK, nil
	}

	if data.DebutDay > 0 {
		newFieldStage = m.addField(newFieldStage, "debut_day", bson.M{"$dayOfMonth": "$debut_date"})
		matchStage = m.addMatch(matchStage, "debut_day", data.DebutDay)
	}

	if data.StartDebutMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "debut_month", bson.M{"$month": "$debut_date"})
		matchStage = m.addMatch(matchStage, "debut_month", bson.M{"$gte": data.StartDebutMonth})
	}

	if data.EndDebutMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "debut_month", bson.M{"$month": "$debut_date"})
		matchStage = m.addMatch(matchStage, "debut_month", bson.M{"$lte": data.EndDebutMonth})
	}

	if data.StartDebutYear > 0 {
		newFieldStage = m.addField(newFieldStage, "debut_year", bson.M{"$year": "$debut_date"})
		matchStage = m.addMatch(matchStage, "debut_year", bson.M{"$gte": data.StartDebutYear})
	}

	if data.EndDebutYear > 0 {
		newFieldStage = m.addField(newFieldStage, "debut_year", bson.M{"$year": "$debut_date"})
		matchStage = m.addMatch(matchStage, "debut_year", bson.M{"$lte": data.EndDebutYear})
	}

	if data.StartRetiredMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "retired_month", bson.M{"$month": "$retirement_date"})
		matchStage = m.addMatch(matchStage, "retired_month", bson.M{"$gte": data.StartRetiredMonth})
	}

	if data.EndRetiredMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "retired_month", bson.M{"$month": "$retirement_date"})
		matchStage = m.addMatch(matchStage, "retired_month", bson.M{"$lte": data.EndRetiredMonth})
	}

	if data.StartRetiredYear > 0 {
		newFieldStage = m.addField(newFieldStage, "retired_year", bson.M{"$year": "$retirement_date"})
		matchStage = m.addMatch(matchStage, "retired_year", bson.M{"$gte": data.StartRetiredYear})
	}

	if data.EndRetiredYear > 0 {
		newFieldStage = m.addField(newFieldStage, "retired_year", bson.M{"$year": "$retirement_date"})
		matchStage = m.addMatch(matchStage, "retired_year", bson.M{"$lte": data.EndRetiredYear})
	}

	if data.Has2D != nil {
		matchStage = m.addMatch(matchStage, "has_2d", utils.PtrToBool(data.Has2D))
	}

	if data.Has3D != nil {
		matchStage = m.addMatch(matchStage, "has_3d", utils.PtrToBool(data.Has3D))
	}

	if data.CharacterDesigner != "" {
		matchStage = m.addMatch(matchStage, "character_designers", data.CharacterDesigner)
	}

	if data.Character2DModeler != "" {
		matchStage = m.addMatch(matchStage, "character_2d_modelers", data.Character2DModeler)
	}

	if data.Character3DModeler != "" {
		matchStage = m.addMatch(matchStage, "character_3d_modelers", data.Character3DModeler)
	}

	if data.InAgency != nil {
		matchStage = m.addMatch(matchStage, "agencies.0", bson.M{"$exists": utils.PtrToBool(data.InAgency)})
	}

	if data.Agency != "" {
		matchStage = m.addMatch(matchStage, "agencies.name", data.Agency)
	}

	if data.AgencyID > 0 {
		matchStage = m.addMatch(matchStage, "agencies.id", data.AgencyID)
	}

	if data.LanguageID > 0 {
		matchStage = m.addMatch(matchStage, "languages.id", data.LanguageID)
	}

	if len(data.ChannelTypes) > 0 {
		matchStage = m.addMatch(matchStage, "channels.type", m.getChannelTypeFilter(data.ChannelTypes))
	}

	if data.BirthdayDay > 0 {
		newFieldStage = m.addField(newFieldStage, "birthday_day", bson.M{"$dayOfMonth": "$birthday"})
		matchStage = m.addMatch(matchStage, "birthday_day", data.BirthdayDay)
	}

	if data.StartBirthdayMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "birthday_month", bson.M{"$month": "$birthday"})
		matchStage = m.addMatch(matchStage, "birthday_month", bson.M{"$gte": data.StartBirthdayMonth})
	}

	if data.EndBirthdayMonth > 0 {
		newFieldStage = m.addField(newFieldStage, "birthday_month", bson.M{"$month": "$birthday"})
		matchStage = m.addMatch(matchStage, "birthday_month", bson.M{"$lte": data.EndBirthdayMonth})
	}

	if len(data.BloodTypes) > 0 {
		matchStage = m.addMatch(matchStage, "blood_type", m.getArrayFilter(data.BloodTypes))
	}

	if len(data.Genders) > 0 {
		matchStage = m.addMatch(matchStage, "gender", m.getArrayFilter(data.Genders))
	}

	if len(data.Zodiacs) > 0 {
		matchStage = m.addMatch(matchStage, "zodiac_sign", m.getArrayFilter(data.Zodiacs))
	}

	if data.StartSubscriber > 0 {
		matchStage = m.addMatch(matchStage, "subscriber", bson.M{"$gte": data.StartSubscriber})
	}

	if data.EndSubscriber > 0 {
		matchStage = m.addMatch(matchStage, "subscriber", bson.M{"$lt": data.EndSubscriber})
	}

	if data.StartVideoCount > 0 {
		matchStage = m.addMatch(matchStage, "video_count", bson.M{"$gte": data.StartVideoCount})
	}

	if data.EndVideoCount > 0 {
		matchStage = m.addMatch(matchStage, "video_count", bson.M{"$lt": data.EndVideoCount})
	}

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, projectStage, sortStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var vtubers []vtuber
	if err := cursor.All(ctx, &vtubers); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]entity.Vtuber, len(vtubers))
	for i, vtuber := range vtubers {
		res[i] = *vtuber.toEntity()
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, countStage))
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

// GetCharacterDesigners to get character designers.
func (m *Mongo) GetCharacterDesigners(ctx context.Context) ([]string, int, error) {
	var designers []string
	if err := m.db.Distinct(ctx, "character_designers", bson.M{"character_designers": bson.M{"$ne": nil}}).Decode(&designers); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return designers, http.StatusOK, nil
}

// GetCharacter2DModelers to get 2d modelers.
func (m *Mongo) GetCharacter2DModelers(ctx context.Context) ([]string, int, error) {
	var modelers []string
	if err := m.db.Distinct(ctx, "character_2d_modelers", bson.M{"character_2d_modelers": bson.M{"$ne": nil}}).Decode(&modelers); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return modelers, http.StatusOK, nil
}

// GetCharacter3DModelers to get 3d modelers.
func (m *Mongo) GetCharacter3DModelers(ctx context.Context) ([]string, int, error) {
	var modelers []string
	if err := m.db.Distinct(ctx, "character_3d_modelers", bson.M{"character_3d_modelers": bson.M{"$ne": nil}}).Decode(&modelers); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return modelers, http.StatusOK, nil
}

// GetVideos to get videos.
func (m *Mongo) GetVideos(ctx context.Context, data entity.GetVideosRequest) ([]entity.VtuberVideo, int, int, error) {
	unwindStage := bson.D{{Key: "$unwind", Value: "$channels"}}
	unwindStage2 := bson.D{{Key: "$unwind", Value: "$channels.videos"}}
	projectStage := bson.D{{Key: "$project", Value: bson.M{
		"id":               "$id",
		"vtuber_id":        "$id",
		"vtuber_name":      "$name",
		"vtuber_image":     "$image",
		"channel_id":       "$channels.id",
		"channel_name":     "$channels.name",
		"channel_type":     "$channels.type",
		"channel_url":      "$channels.url",
		"video_id":         "$channels.videos.id",
		"video_title":      "$channels.videos.title",
		"video_url":        "$channels.videos.url",
		"video_image":      "$channels.videos.image",
		"video_start_date": "$channels.videos.start_date",
		"video_end_date":   "$channels.videos.end_date",
	}}}
	matchStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: m.convertSort(data.Sort)}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.StartDate != nil {
		matchStage = m.addMatch(matchStage, "video_start_date", bson.M{"$gte": bson.NewDateTimeFromTime(*data.StartDate)})
	}

	if data.EndDate != nil {
		matchStage = m.addMatch(matchStage, "video_start_date", bson.M{"$lte": bson.NewDateTimeFromTime(*data.EndDate)})
	}

	if data.IsFinished != nil {
		key := map[bool]string{false: "$eq", true: "$ne"}
		matchStage = m.addMatch(matchStage, "video_end_date", bson.M{key[*data.IsFinished]: nil})
	}

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(unwindStage, unwindStage2, projectStage, matchStage, sortStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var videos []vtuberVideo
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]entity.VtuberVideo, len(videos))
	for i, video := range videos {
		res[i] = entity.VtuberVideo{
			VtuberID:       video.VtuberID,
			VtuberName:     video.VtuberName,
			VtuberImage:    video.VtuberImage,
			ChannelID:      video.ChannelID,
			ChannelName:    video.ChannelName,
			ChannelType:    video.ChannelType,
			ChannelURL:     video.ChannelURL,
			VideoID:        video.VideoID,
			VideoTitle:     video.VideoTitle,
			VideoURL:       video.VideoURL,
			VideoImage:     video.VideoImage,
			VideoStartDate: video.VideoStartDate,
			VideoEndDate:   video.VideoEndDate,
		}
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(unwindStage, unwindStage2, projectStage, matchStage, countStage))
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

// UpdateOverriddenFieldByID to update overriden field by id.
func (m *Mongo) UpdateOverriddenFieldByID(ctx context.Context, id int64, data entity.OverriddenField) (int, error) {
	if _, err := m.db.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{
		"overridden_field": m.overiddenFieldFromEntity(data),
		"updated_at":       time.Now(),
	}}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}
