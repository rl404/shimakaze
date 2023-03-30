package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
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

// DeleteByID to delete by id.
func (m *Mongo) DeleteByID(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
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
	newFieldStage := bson.D{}
	matchStage := bson.D{}
	omitStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: m.convertSort(data.Sort)}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.Mode == entity.SearchModeStats {
		omitStage = append(omitStage, bson.E{Key: "$unset", Value: bson.A{
			"original_names",
			"nicknames",
			"caption",
			"affiliations",
			"official_websites",
		}})
	} else {
		omitStage = append(omitStage, bson.E{Key: "$unset", Value: bson.A{
			"channels.videos",
		}})
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

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, omitStage, sortStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var vtubers []vtuber
	if err := cursor.All(ctx, &vtubers); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]entity.Vtuber, len(vtubers))
	for i, vtuber := range vtubers {
		res[i] = *vtuber.toEntity()
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, countStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total []map[string]int64
	if err := cntCursor.All(ctx, &total); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if len(total) == 0 {
		return res, 0, http.StatusOK, nil
	}

	return res, int(total[0]["count"]), http.StatusOK, nil
}

// GetCharacterDesigners to get character designers.
func (m *Mongo) GetCharacterDesigners(ctx context.Context) ([]string, int, error) {
	designers, err := m.db.Distinct(ctx, "character_designers", bson.M{"character_designers": bson.M{"$ne": nil}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]string, len(designers))
	for i, d := range designers {
		v, ok := d.(string)
		if !ok {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, _errors.New("invalid value"))
		}
		res[i] = v
	}

	return res, http.StatusOK, nil
}

// GetCharacter2DModelers to get 2d modelers.
func (m *Mongo) GetCharacter2DModelers(ctx context.Context) ([]string, int, error) {
	modelers, err := m.db.Distinct(ctx, "character_2d_modelers", bson.M{"character_2d_modelers": bson.M{"$ne": nil}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]string, len(modelers))
	for i, d := range modelers {
		v, ok := d.(string)
		if !ok {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, _errors.New("invalid value"))
		}
		res[i] = v
	}

	return res, http.StatusOK, nil
}

// GetCharacter3DModelers to get 3d modelers.
func (m *Mongo) GetCharacter3DModelers(ctx context.Context) ([]string, int, error) {
	modelers, err := m.db.Distinct(ctx, "character_3d_modelers", bson.M{"character_3d_modelers": bson.M{"$ne": nil}})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]string, len(modelers))
	for i, d := range modelers {
		v, ok := d.(string)
		if !ok {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, _errors.New("invalid value"))
		}
		res[i] = v
	}

	return res, http.StatusOK, nil
}
