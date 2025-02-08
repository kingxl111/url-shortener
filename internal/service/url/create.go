package url

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/service/shortener"

	"github.com/kingxl111/url-shortener/internal/model"
)

func (s *service) CreateURL(ctx context.Context, url model.URL) (model.URL, error) {
	// TODO: errors handling
	url.ShortenedURL = shortener.GenerateShortURL(url.OriginalURL)
	return s.urlRepository.Create(ctx, url)
}
