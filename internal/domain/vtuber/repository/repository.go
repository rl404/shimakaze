package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

// Repository contains functions for vtuber domain.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*entity.Vtuber, int, error)
	UpdateByID(ctx context.Context, id int64, data entity.Vtuber) (int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)
	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetOldIDs(ctx context.Context) ([]int64, int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Vtuber, int, int, error)
	GetAllIDs(ctx context.Context) ([]int64, int, error)
	GetAllImages(ctx context.Context, shuffle bool, limit int) ([]entity.Vtuber, int, error)
	GetAllForFamilyTree(ctx context.Context) ([]entity.Vtuber, int, error)
	GetAllForAgencyTree(ctx context.Context) ([]entity.Vtuber, int, error)
	GetCharacterDesigners(ctx context.Context) ([]string, int, error)
	GetCharacter2DModelers(ctx context.Context) ([]string, int, error)
	GetCharacter3DModelers(ctx context.Context) ([]string, int, error)
	GetCount(ctx context.Context) (int, int, error)
	GetAverageActiveTime(ctx context.Context) (float64, int, error)
	GetStatusCount(ctx context.Context) (*entity.StatusCount, int, error)
	GetDebutRetireCountMonthly(ctx context.Context) ([]entity.DebutRetireCount, int, error)
	GetDebutRetireCountYearly(ctx context.Context) ([]entity.DebutRetireCount, int, error)
}
