package grpc

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/url"
)

type (
	ShortenerService interface {
		CreateURL(ctx context.Context, url url.URL) (*url.URL, error)
		GetURL(ctx context.Context, shortenedUrl url.URL) (*url.URL, error)
	}
)
