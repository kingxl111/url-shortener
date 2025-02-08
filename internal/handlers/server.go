package handlers

import (
	"context"
	"fmt"
	"github.com/kingxl111/url-shortener/internal/model"
	"log"

	serv "github.com/kingxl111/url-shortener/internal/service"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
)

type Server struct {
	shrt.UnimplementedURLShortenerServer
	Services serv.ShortenerService
}

func (s *Server) Create(ctx context.Context, req *shrt.CreateRequest) (*shrt.CreateResponse, error) {
	url := model.URL{OriginalURL: req.OriginalUrl}
	shortenedURL, err := s.Services.CreateURL(ctx, url)
	if err != nil {
		log.Printf("error in create: %v", err)
		return nil, fmt.Errorf("create shortened url: %w", err)
	}

	return &shrt.CreateResponse{
		ShortUrl: shortenedURL.ShortenedURL,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *shrt.GetRequest) (*shrt.GetResponse, error) {
	url := model.URL{ShortenedURL: req.ShortUrl}
	originalURL, err := s.Services.GetURL(ctx, url)
	if err != nil {
		log.Printf("error in getting url: %v", err)
		return nil, fmt.Errorf("get original url: %w", err)
	}

	return &shrt.GetResponse{
		OriginalUrl: originalURL.OriginalURL,
	}, nil
}
