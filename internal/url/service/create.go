package service

import (
	"context"
	"errors"
	"fmt"
	net "net/url"
	"strings"

	"github.com/kingxl111/url-shortener/internal/repository"
	urlPack "github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
)

func ValidateAndNormalizeURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", urlPack.ErrEmptyURL
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsed, err := net.ParseRequestURI(rawURL)
	if err != nil {
		return "", urlPack.ErrInvalidFormat
	}

	if strings.ContainsAny(parsed.Hostname(), " !{}|\\^\"<>") {
		return "", urlPack.ErrInvalidFormat
	}

	if parsed.Host == "" {
		return "", urlPack.ErrMissingHost
	}

	if strings.Contains(parsed.Host, " ") {
		return "", urlPack.ErrInvalidFormat
	}

	allowedSchemes := map[string]bool{"http": true, "https": true}
	if !allowedSchemes[parsed.Scheme] {
		return "", urlPack.ErrInvalidScheme
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)

	return parsed.String(), nil
}

func (s *service) CreateURL(ctx context.Context, inputURL urlPack.URL) (*urlPack.URL, error) {
	normalized, err := ValidateAndNormalizeURL(inputURL.OriginalURL)
	if err != nil {
		return nil, err
	}

	inputURL.OriginalURL = normalized
	inputURL.ShortenedURL = shortener.GenerateShortURL(inputURL.OriginalURL)
	shortenedURL, err := s.urlRepository.Create(ctx, inputURL)
	if err != nil {
		if errors.Is(err, repository.ErrorDuplicatedURL) {
			return nil, fmt.Errorf("url already exists")
		}
		return nil, fmt.Errorf("internal server error")
	}

	return shortenedURL, nil
}
