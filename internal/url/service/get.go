package service

import (
	"context"
	"errors"

	"github.com/kingxl111/url-shortener/internal/repository"

	urlPack "github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
)

func (s *service) GetURL(ctx context.Context, shortenedURL urlPack.URL) (*urlPack.URL, error) {
	if shortenedURL.ShortenedURL == "" {
		return nil, urlPack.ErrEmptyURL
	}

	if len(shortenedURL.ShortenedURL) != shortener.ShortURLLength {
		return nil, urlPack.ErrInvalidLength
	}

	if !IsValidShortURL(shortenedURL.ShortenedURL) {
		return nil, urlPack.ErrInvalidCharacters
	}

	result, err := s.urlRepository.Get(ctx, shortenedURL)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return nil, repository.ErrorNotFound
		}
		return nil, urlPack.ErrRepository
	}

	return result, nil
}

func IsValidShortURL(s string) bool {
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
