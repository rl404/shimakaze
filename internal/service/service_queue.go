package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/language/entity"
)

// QueueMissingVtuber to queue missing vtuber.
func (s *service) QueueMissingVtuber(ctx context.Context, limit int) (int, int, error) {
	vtuberIDs, code, err := s.vtuber.GetAllIDs(ctx)
	if err != nil {
		return 0, code, stack.Wrap(ctx, err)
	}

	nonVtuberIDs, code, err := s.nonVtuber.GetAllIDs(ctx)
	if err != nil {
		return 0, code, stack.Wrap(ctx, err)
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
			return cnt, code, stack.Wrap(ctx, err)
		}

		lastName = nextName

		for _, page := range pages {
			if existMap[page.ID] {
				continue
			}

			existMap[page.ID] = true

			if err := s.publisher.PublishParseVtuber(ctx, page.ID, false); err != nil {
				return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
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
func (s *service) QueueMissingAgency(ctx context.Context, limit int) (int, int, error) {
	agencyIDs, code, err := s.agency.GetAllIDs(ctx)
	if err != nil {
		return 0, code, stack.Wrap(ctx, err)
	}

	existMap := make(map[int64]bool)
	for _, id := range agencyIDs {
		existMap[id] = true
	}

	var cnt int
	var lastTitle string
	limitPerPage := 500
	for {
		agencies, nextTitle, code, err := s.wikia.GetCategoryMembers(ctx, "Category:Agency", limitPerPage, lastTitle, true)
		if err != nil {
			return cnt, code, stack.Wrap(ctx, err)
		}

		lastTitle = nextTitle

		for _, agency := range agencies {
			if existMap[agency.ID] {
				continue
			}

			existMap[agency.ID] = true

			if err := s.publisher.PublishParseAgency(ctx, agency.ID, false); err != nil {
				return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
			}

			cnt++

			if cnt >= limit {
				return cnt, http.StatusOK, nil
			}
		}

		if len(agencies) == 0 || lastTitle == "" {
			return cnt, http.StatusOK, nil
		}
	}
}

// QueueMissingLanguage to queue missing language.
func (s *service) QueueMissingLanguage(ctx context.Context, limit int) (int, int, error) {
	languageIDs, code, err := s.language.GetAllIDs(ctx)
	if err != nil {
		return 0, code, stack.Wrap(ctx, err)
	}

	existMap := make(map[int64]bool)
	for _, id := range languageIDs {
		existMap[id] = true
	}

	var cnt int
	var lastTitle string
	limitPerPage := 500
	for {
		languages, nextTitle, code, err := s.wikia.GetCategoryMembers(ctx, "Category:Language", limitPerPage, lastTitle, false)
		if err != nil {
			return cnt, code, stack.Wrap(ctx, err)
		}

		lastTitle = nextTitle

		for _, language := range languages {
			if existMap[language.ID] {
				continue
			}

			existMap[language.ID] = true

			if code, err := s.language.UpdateByID(ctx, language.ID, entity.Language{
				ID:   language.ID,
				Name: strings.ReplaceAll(language.Title, "Category:", ""),
			}); err != nil {
				return cnt, code, stack.Wrap(ctx, err)
			}

			cnt++

			if cnt >= limit {
				return cnt, http.StatusOK, nil
			}
		}

		if len(languages) == 0 || lastTitle == "" {
			return cnt, http.StatusOK, nil
		}
	}
}

// QueueOldAgency to queue old agency.
func (s *service) QueueOldAgency(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.agency.GetOldIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for _, id := range ids {
		if err := s.publisher.PublishParseAgency(ctx, id, false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}

		cnt++

		if cnt >= limit {
			break
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldActiveVtuber to queue old active vtuber.
func (s *service) QueueOldActiveVtuber(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.vtuber.GetOldActiveIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for _, id := range ids {
		if err := s.publisher.PublishParseVtuber(ctx, id, false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}

		cnt++

		if cnt >= limit {
			break
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldRetiredVtuber to queue old retired vtuber.
func (s *service) QueueOldRetiredVtuber(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.vtuber.GetOldRetiredIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for _, id := range ids {
		if err := s.publisher.PublishParseVtuber(ctx, id, false); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}

		cnt++

		if cnt >= limit {
			break
		}
	}

	return cnt, http.StatusOK, nil
}
