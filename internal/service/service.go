package service

import (
	"context"

	agencyRepository "github.com/rl404/shimakaze/internal/domain/agency/repository"
	nonVtuberRepository "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	publisherRepository "github.com/rl404/shimakaze/internal/domain/publisher/repository"
	vtuberRepository "github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	wikiaRepository "github.com/rl404/shimakaze/internal/domain/wikia/repository"
)

// Service contains functions for service.
type Service interface {
	GetVtuberByID(ctx context.Context, id int64) (*vtuber, int, error)
	GetVtuberImages(ctx context.Context) ([]vtuberImage, int, error)
	GetVtuberFamilyTrees(ctx context.Context) (*vtuberFamilyTree, int, error)

	GetWikiaImage(ctx context.Context, path string) ([]byte, int, error)

	ConsumeMessage(ctx context.Context, msg entity.Message) error

	QueueMissingAgency(ctx context.Context) (int, int, error)
	QueueMissingVtuber(ctx context.Context) (int, int, error)
	QueueOldAgency(ctx context.Context) (int, int, error)
	QueueOldVtuber(ctx context.Context) (int, int, error)
}

type service struct {
	wikia     wikiaRepository.Repository
	vtuber    vtuberRepository.Repository
	nonVtuber nonVtuberRepository.Repository
	agency    agencyRepository.Repository
	publisher publisherRepository.Repository
}

// New to create new service.
func New(
	wikia wikiaRepository.Repository,
	vtuber vtuberRepository.Repository,
	nonVtuber nonVtuberRepository.Repository,
	agency agencyRepository.Repository,
	publisher publisherRepository.Repository,
) Service {
	return &service{
		wikia:     wikia,
		vtuber:    vtuber,
		nonVtuber: nonVtuber,
		agency:    agency,
		publisher: publisher,
	}
}
