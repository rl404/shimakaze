package mongo

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/user/entity"
	"github.com/rl404/shimakaze/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo contains function for user mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new user mongodb.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("user"),
	}
}

// Upsert to upsert user.
func (m *Mongo) Upsert(ctx context.Context, data entity.User) (int, error) {
	var user user
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&user); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := m.db.InsertOne(ctx, m.fromEntity(data)); err != nil {
				return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	uu := m.fromEntity(data)
	uu.CreatedAt = user.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": uu}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.User, int, error) {
	var user user
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrUserNotFound)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return user.toEntity(), http.StatusOK, nil
}
