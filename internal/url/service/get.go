package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
)

func (s *service) GetURL(ctx context.Context, shortenedURL url.URL) (*url.URL, error) {
	if shortenedURL.ShortenedURL == "" {
		return nil, fmt.Errorf("shortened URL is empty")
	}

	if len(shortenedURL.ShortenedURL) != shortener.ShortURLLength {
		return nil, url.ErrInvalidLength
	}

	if !isValidShortURL(shortenedURL.ShortenedURL) {
		return nil, url.ErrInvalidCharacters
	}

	result, err := s.urlRepository.Get(ctx, shortenedURL)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return nil, repository.ErrorNotFound
		}
		return nil, url.ErrRepository
	}

	return result, nil
}

func isValidShortURL(s string) bool {
	for _, c := range s {
		if !isAllowedShortURLChar(c) {
			return false
		}
	}
	return true
}

func isAllowedShortURLChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
