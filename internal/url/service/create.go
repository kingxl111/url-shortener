package service

import (
	"context"
	"fmt"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
	net "net/url"
	"strings"
)

func ValidateAndNormalizeURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("URL is empty")
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsed, err := net.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("missing host in URL")
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)

	return parsed.String(), nil
}

func (s *service) CreateURL(ctx context.Context, inputURL url.URL) (*url.URL, error) {
	normalized, err := ValidateAndNormalizeURL(inputURL.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	inputURL.OriginalURL = normalized
	inputURL.ShortenedURL = shortener.GenerateShortURL(inputURL.OriginalURL)
	shortenedURL, err := s.urlRepository.Create(ctx, inputURL)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return shortenedURL, nil
}
