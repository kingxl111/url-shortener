package url

import (
	"context"
)

type ShortenerService interface {
	CreateURL(ctx context.Context, url URL) (*URL, error)
	GetURL(ctx context.Context, shortenedUrl URL) (*URL, error)
}
