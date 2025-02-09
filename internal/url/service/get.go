package service

import (
	"context"
	"errors"

	"github.com/kingxl111/url-shortener/internal/repository"

	"github.com/kingxl111/url-shortener/internal/url"
)

func (s *service) GetURL(ctx context.Context, shortenedURL url.URL) (*url.URL, error) {
	result, err := s.urlRepository.Get(ctx, shortenedURL)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return nil, url.ErrNotFoundURL
		}
		return nil, url.ErrRepository
	}

	return result, nil
}
