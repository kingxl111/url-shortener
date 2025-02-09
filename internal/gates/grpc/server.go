package grpc

import (
	"context"
	"fmt"

	urlPck "github.com/kingxl111/url-shortener/internal/url"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
)

type Server struct {
	shrt.UnimplementedURLShortenerServer
	Services urlPck.ShortenerService
}

func (s *Server) Create(ctx context.Context, req *shrt.Create_Request) (*shrt.Create_Response, error) {
	url := urlPck.URL{OriginalURL: req.OriginalUrl}
	shortenedURL, err := s.Services.CreateURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("create shortened url: %w", err)
	}

	return &shrt.Create_Response{
		ShortUrl: shortenedURL.ShortenedURL,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *shrt.Get_Request) (*shrt.Get_Response, error) {
	url := urlPck.URL{ShortenedURL: req.ShortUrl}
	originalURL, err := s.Services.GetURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("get original url: %w", err)
	}

	return &shrt.Get_Response{
		OriginalUrl: originalURL.OriginalURL,
	}, nil
}
