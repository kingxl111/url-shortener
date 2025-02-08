package repository

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/url"
)

type URLRepository interface {
	Create(ctx context.Context, url url.URL) (*url.URL, error)
	Get(ctx context.Context, shortenedUrl url.URL) (*url.URL, error)
}
