package repository

import (
	"context"

	"github.com/kingxl111/url-shortener/internal/url"
)

//go:generate mockgen -source=repository.go -destination=mocks/mocks.go -package=repository
type URLRepository interface {
	Create(ctx context.Context, url url.URL) (*url.URL, error)
	Get(ctx context.Context, shortenedUrl url.URL) (*url.URL, error)
}
