package service

import (
	"context"
	"net/http"

	"github.com/rl404/shimakaze/internal/errors"
)

// GetWikiaImage to get wikia image.
func (s *service) GetWikiaImage(ctx context.Context, path string) ([]byte, int, error) {
	image, code, err := s.wikia.GetImage(ctx, path)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}
	return image, http.StatusOK, nil
}
