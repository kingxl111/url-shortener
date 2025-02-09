package service

import (
	"context"
	"errors"
	"github.com/kingxl111/url-shortener/internal/repository"

	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
)

func (s *service) CreateURL(ctx context.Context, inputURL url.URL) (*url.URL, error) {

	inputURL.ShortenedURL = shortener.GenerateShortURL(inputURL.OriginalURL)
	shortenedURL, err := s.urlRepository.Create(ctx, inputURL)
	if err != nil {
		if errors.Is(err, repository.ErrorDuplicatedURL) {
			return nil, url.ErrDuplicatedURL
		}
		return nil, url.ErrService
	}

	return shortenedURL, nil
}
