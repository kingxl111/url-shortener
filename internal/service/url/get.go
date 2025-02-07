package url

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/model"
)

func (s *service) GetURL(ctx context.Context, shortenedUrl model.URL) (model.URL, error) {
	return s.GetURL(ctx, shortenedUrl)
}
