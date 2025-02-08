package url

import (
	"context"

	"github.com/kingxl111/url-shortener/internal/model"
)

func (s *service) CreateURL(ctx context.Context, url model.URL) (model.URL, error) {
	url.ShortenedURL = "beer"
	return s.urlRepository.Create(ctx, url)
}
