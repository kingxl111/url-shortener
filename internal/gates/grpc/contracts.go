package grpc

import (
	"context"

	"github.com/kingxl111/url-shortener/internal/url"
)

//go:generate mockgen -source=contracts.go -destination=mocks.go -package=grpc
type (
	ShortenerService interface {
		CreateURL(ctx context.Context, inputURL url.URL) (*url.URL, error)
		GetURL(ctx context.Context, shortenedUrl url.URL) (*url.URL, error)
	}
)
