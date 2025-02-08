package service

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/url"
)

func (s *service) GetURL(ctx context.Context, shortenedUrl url.URL) (*url.URL, error) {
	return s.urlRepository.Get(ctx, shortenedUrl)
}
