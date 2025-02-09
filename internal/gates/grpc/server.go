package grpc

import (
	"context"
	"errors"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	net "net/url"
	"strings"

	"github.com/kingxl111/url-shortener/internal/url"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
)

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

	if !IsValidShortURL(req.ShortUrl) {
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

func IsValidShortURL(s string) bool {
	for _, c := range s {
		if !isAllowedShortURLChar(c) {
			return false
		}
	}
	return true
}

func isAllowedShortURLChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
