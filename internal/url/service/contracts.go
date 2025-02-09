package service

import (
	"context"

	"github.com/kingxl111/url-shortener/internal/url"
)

//go:generate mockgen -source=contracts.go -destination=mocks.go -package=service
type (
	URLRepository interface {
		Create(ctx context.Context, inputURL url.URL) (*url.URL, error)
		Get(ctx context.Context, shortenedUrl url.URL) (*url.URL, error)
	}
)
