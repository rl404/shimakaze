package service

import (
	"context"
	"net/http"
	"net/url"

	"github.com/rl404/fairy/errors/stack"
)

// GetWikiaImage to get wikia image.
func (s *service) GetWikiaImage(ctx context.Context, path string) ([]byte, int, error) {
	path, _ = url.QueryUnescape(path)
	image, code, err := s.wikia.GetImage(ctx, path)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return image, http.StatusOK, nil
}
