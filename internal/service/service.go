package service

import (
	"context"

	"github.com/kingxl111/url-shortener/internal/model"
)

type ShortenerService interface {
	CreateURL(ctx context.Context, url model.URL) (model.URL, error)
	GetURL(ctx context.Context, shortenedUrl model.URL) (model.URL, error)
}
