package service

import (
	"context"

	agencyRepository "github.com/rl404/shimakaze/internal/domain/agency/repository"
	bilibilRepository "github.com/rl404/shimakaze/internal/domain/bilibili/repository"
	niconicoRepository "github.com/rl404/shimakaze/internal/domain/niconico/repository"
	nonVtuberRepository "github.com/rl404/shimakaze/internal/domain/non_vtuber/repository"
	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	publisherRepository "github.com/rl404/shimakaze/internal/domain/publisher/repository"
	twitchRepository "github.com/rl404/shimakaze/internal/domain/twitch/repository"
	vtuberRepository "github.com/rl404/shimakaze/internal/domain/vtuber/repository"
	wikiaRepository "github.com/rl404/shimakaze/internal/domain/wikia/repository"
	youtubeRepository "github.com/rl404/shimakaze/internal/domain/youtube/repository"
)

// Service contains functions for service.
type Service interface {
	GetVtubers(ctx context.Context, params GetVtubersRequest) ([]vtuber, *pagination, int, error)
	GetVtuberByID(ctx context.Context, id int64) (*vtuber, int, error)
	GetVtuberImages(ctx context.Context, shuffle bool, limit int) ([]vtuberImage, int, error)
	GetVtuberFamilyTrees(ctx context.Context) (*vtuberFamilyTree, int, error)
	GetVtuberAgencyTrees(ctx context.Context) (*vtuberAgencyTree, int, error)
	GetVtuberCharacterDesigners(ctx context.Context) ([]string, int, error)
	GetVtuberCharacter2DModelers(ctx context.Context) ([]string, int, error)
	GetVtuberCharacter3DModelers(ctx context.Context) ([]string, int, error)
	GetVtuberCount(ctx context.Context) (int, int, error)
	GetVtuberAverageActiveTime(ctx context.Context) (float64, int, error)
	GetVtuberStatusCount(ctx context.Context) (*vtuberStatusCount, int, error)
	GetVtuberDebutRetireCountMonthly(ctx context.Context) ([]vtuberDebutRetireCount, int, error)
	GetVtuberDebutRetireCountYearly(ctx context.Context) ([]vtuberDebutRetireCount, int, error)
	GetVtuberModelCount(ctx context.Context) (*vtuberModelCount, int, error)
	GetVtuberInAgencyCount(ctx context.Context) (*vtuberInAgencyCount, int, error)
	GetVtuberSubscriberCount(ctx context.Context, params GetVtuberSubscriberCountRequest) ([]vtuberSubscriberCount, int, error)
	GetVtuberDesignerCount(ctx context.Context, params GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error)
	GetVtuber2DModelerCount(ctx context.Context, params GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error)
	GetVtuber3DModelerCount(ctx context.Context, params GetVtuberDesignerCountRequest) ([]vtuberDesignerCount, int, error)

	GetAgencies(ctx context.Context, params GetAgenciesRequest) ([]agency, *pagination, int, error)
	GetAgencyByID(ctx context.Context, id int64) (*agency, int, error)
	GetAgencyCount(ctx context.Context) (int, int, error)

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
	youtube   youtubeRepository.Repository
	twitch    twitchRepository.Repository
	bilibili  bilibilRepository.Repository
	niconico  niconicoRepository.Repository
}

// New to create new service.
func New(
	wikia wikiaRepository.Repository,
	vtuber vtuberRepository.Repository,
	nonVtuber nonVtuberRepository.Repository,
	agency agencyRepository.Repository,
	publisher publisherRepository.Repository,
	youtube youtubeRepository.Repository,
	twitch twitchRepository.Repository,
	bilibili bilibilRepository.Repository,
	niconico niconicoRepository.Repository,
) Service {
	return &service{
		wikia:     wikia,
		vtuber:    vtuber,
		nonVtuber: nonVtuber,
		agency:    agency,
		publisher: publisher,
		youtube:   youtube,
		twitch:    twitch,
		bilibili:  bilibili,
		niconico:  niconico,
	}
}

type pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
