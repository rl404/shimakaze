package repository

import (
	"context"

	"github.com/rl404/shimakaze/internal/domain/wikia/entity"
)

// Repository contains functions for wikia domain.
type Repository interface {
	GetPages(ctx context.Context, apLimit int, apContinue string) ([]entity.Page, string, int, error)
	GetPageByID(ctx context.Context, id int64) (*entity.Page, int, error)
	GetPageImageByID(ctx context.Context, id int64) (*entity.PageImage, int, error)
	GetCategoryMembers(ctx context.Context, cmTitle string, cmLimit int, cmContinue string) ([]entity.CategoryMember, string, int, error)
	GetImageInfo(ctx context.Context, imageName string) (string, int, error)
	GetPageCategories(ctx context.Context, id int64, clLimit int, clContinue string) ([]entity.PageCategory, string, int, error)

	GetImage(ctx context.Context, path string) ([]byte, int, error)
}
