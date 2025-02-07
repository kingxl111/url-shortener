package repository

import (
	"context"
	"github.com/kingxl111/url-shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url model.URL) (model.URL, error)
	Get(ctx context.Context, shortenedUrl model.URL) (model.URL, error)
}
