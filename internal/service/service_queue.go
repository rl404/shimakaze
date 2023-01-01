package service

import (
	"context"
	"net/http"

	"github.com/rl404/shimakaze/internal/domain/publisher/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

// QueueMissingVtuber to queue missing vtuber.
func (s *service) QueueMissingVtuber(ctx context.Context, limit int) (int, int, error) {
	vtuberIDs, code, err := s.vtuber.GetAllIDs(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	nonVtuberIDs, code, err := s.nonVtuber.GetAllIDs(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	existMap := make(map[int64]bool)
	for _, id := range append(vtuberIDs, nonVtuberIDs...) {
		existMap[id] = true
	}

	var cnt int
	var lastName string
	limitPerPage := 500
	for {
		pages, nextName, code, err := s.wikia.GetPages(ctx, limitPerPage, lastName)
		if err != nil {
			return cnt, code, errors.Wrap(ctx, err)
		}

		lastName = nextName

		for _, page := range pages {
			if existMap[page.ID] {
				continue
			}

			existMap[page.ID] = true

			if err := s.publisher.PublishParseVtuber(ctx, entity.ParseVtuberRequest{ID: page.ID}); err != nil {
				return cnt, http.StatusInternalServerError, errors.Wrap(ctx, err)
			}

			cnt++
			if cnt >= limit {
				return cnt, http.StatusOK, nil
			}
		}

		if len(pages) == 0 || lastName == "" {
			return cnt, http.StatusOK, nil
		}
	}
}

// QueueMissingAgency to queue missing agency.
func (s *service) QueueMissingAgency(ctx context.Context) (int, int, error) {
	agencyIDs, code, err := s.agency.GetAllIDs(ctx)
	if err != nil {
		return 0, code, errors.Wrap(ctx, err)
	}

	existMap := make(map[int64]bool)
	for _, id := range agencyIDs {
		existMap[id] = true
	}

	var cnt int
	var lastTitle string
	limitPerPage := 500
	for {
		agencies, nextTitle, code, err := s.wikia.GetCategoryMembers(ctx, "Category:Agency", limitPerPage, lastTitle)
		if err != nil {
			return cnt, code, errors.Wrap(ctx, err)
		}

		lastTitle = nextTitle

		for _, agency := range agencies {
			if existMap[agency.ID] {
				continue
			}

			existMap[agency.ID] = true

			if err := s.publisher.PublishParseAgency(ctx, entity.ParseAgencyRequest{ID: agency.ID}); err != nil {
				return cnt, http.StatusInternalServerError, errors.Wrap(ctx, err, errors.ErrInternalServer)
			}

			cnt++
		}

		if len(agencies) == 0 || lastTitle == "" {
			return cnt, http.StatusOK, nil
		}
	}
}
