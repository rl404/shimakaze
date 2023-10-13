package service

import (
	"context"
	"net/http"
	"regexp"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	vtuberEntity "github.com/rl404/shimakaze/internal/domain/vtuber/entity"
)

func (s *service) updateAgency(ctx context.Context, id int64) (int, error) {
	// Call wikia api.
	page, code, err := s.wikia.GetPageByID(ctx, id)
	if err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Get members.
	vtubers, total, code, err := s.vtuber.GetAll(ctx, vtuberEntity.GetAllRequest{
		Mode:     vtuberEntity.SearchModeAll,
		AgencyID: page.ID,
		Page:     1,
		Limit:    -1,
	})
	if err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Get total subs.
	subsTotal, channelMap := 0, make(map[string]bool)
	for _, vtuber := range vtubers {
		max, channelID := 0, ""
		for _, channel := range vtuber.Channels {
			if channel.Subscriber > max {
				max, channelID = channel.Subscriber, channel.ID
			}
		}

		if !channelMap[channelID] {
			subsTotal += max
		}

		channelMap[channelID] = true
	}

	// Update data.
	if code, err := s.agency.UpdateByID(ctx, id, entity.Agency{
		ID:         page.ID,
		Name:       page.Title,
		Image:      s.getAgencyLogo(ctx, page.Content),
		Member:     total,
		Subscriber: subsTotal,
	}); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}

func (s *service) getAgencyLogo(ctx context.Context, data string) string {
	logoRegex := regexp.MustCompile(`\[\[(File:.+?)(\|.+)?\]\]`)
	if logoRegex.FindString(data) == "" {
		return ""
	}

	submatch := logoRegex.FindStringSubmatch(data)

	if len(submatch) < 2 {
		return ""
	}

	imageURL, _, err := s.wikia.GetImageInfo(ctx, submatch[1])
	if err != nil {
		stack.Wrap(ctx, err)
		return ""
	}

	return imageURL
}
