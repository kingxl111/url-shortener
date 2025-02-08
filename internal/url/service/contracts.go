package service

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/url"
)

//go:generate mockgen -source=contracts.go -destination=mocks.go -package=service
type (
	storage interface {
		Create(ctx context.Context, url url.URL) (url.URL, error)
		Get(ctx context.Context, shortenedUrl url.URL) (url.URL, error)
	}
)
