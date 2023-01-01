package service

import (
	"context"
	"net/http"
	"regexp"

	"github.com/rl404/shimakaze/internal/domain/agency/entity"
	"github.com/rl404/shimakaze/internal/errors"
)

func (s *service) updateAgency(ctx context.Context, id int64) (int, error) {
	// Call wikia api.
	page, code, err := s.wikia.GetPageByID(ctx, id)
	if err != nil {
		return code, errors.Wrap(ctx, err)
	}

	// Update data.
	if code, err := s.agency.UpdateByID(ctx, id, entity.Agency{
		ID:    page.ID,
		Name:  page.Title,
		Image: s.getAgencyLogo(ctx, page.Content),
	}); err != nil {
		return code, errors.Wrap(ctx, err)
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
		errors.Wrap(ctx, err)
		return ""
	}

	return imageURL
}
