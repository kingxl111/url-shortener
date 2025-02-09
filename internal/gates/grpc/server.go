package grpc

import (
	"context"
	"errors"
	net "net/url"
	"strings"

	"github.com/kingxl111/url-shortener/internal/url/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kingxl111/url-shortener/internal/url"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
)

const maxUrlSize = 1024

type Server struct {
	shrt.UnimplementedURLShortenerServer
	Services ShortenerService
}

func (s *Server) Create(ctx context.Context, req *shrt.Create_Request) (*shrt.Create_Response, error) {
	err := ValidateURL(req.OriginalUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	u := url.URL{OriginalURL: req.OriginalUrl}
	shortenedURL, err := s.Services.CreateURL(ctx, u)
	if err != nil {
		if errors.Is(err, url.ErrDuplicatedURL) {
			return nil, status.Error(codes.AlreadyExists, "already exists")
		}
		return nil, status.Error(codes.Internal, "error creating shortened url")
	}

	return &shrt.Create_Response{
		ShortUrl: shortenedURL.ShortenedURL,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *shrt.Get_Request) (*shrt.Get_Response, error) {
	if req.ShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "empty short url")
	}

	if len(req.ShortUrl) != shortener.ShortURLLength {
		return nil, status.Error(codes.InvalidArgument, "invalid short url length")
	}

	if !shortener.IsValidShortURL(req.ShortUrl) {
		return nil, status.Error(codes.InvalidArgument, "invalid short url")
	}

	u := url.URL{ShortenedURL: req.ShortUrl}
	originalURL, err := s.Services.GetURL(ctx, u)
	if err != nil {
		if errors.Is(err, url.ErrNotFoundURL) {
			return nil, status.Error(codes.NotFound, "short url not found")
		}
		return nil, status.Error(codes.Internal, "error getting shortened url")
	}

	return &shrt.Get_Response{
		OriginalUrl: originalURL.OriginalURL,
	}, nil
}

func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return url.ErrEmptyURL
	}

	if len(rawURL) > maxUrlSize {
		return url.ErrInvalidFormat
	}

	parsed, err := net.ParseRequestURI(rawURL)
	if err != nil {
		return url.ErrInvalidFormat
	}

	if strings.ContainsAny(parsed.Hostname(), " !{}|\\^\"<>") {
		return url.ErrInvalidFormat
	}

	if parsed.Host == "" {
		return url.ErrMissingHost
	}

	if strings.Contains(parsed.Host, " ") {
		return url.ErrInvalidFormat
	}

	return nil
}
