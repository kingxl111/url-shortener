package grpc

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kingxl111/url-shortener/internal/url"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: "https://example.com"}
	expectedShortURL := "abc1231231"

	mockService.EXPECT().
		CreateURL(ctx, url.URL{OriginalURL: req.OriginalUrl}).
		Return(&url.URL{OriginalURL: req.OriginalUrl, ShortenedURL: expectedShortURL}, nil).
		Times(1)

	resp, err := server.Create(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expectedShortURL, resp.ShortUrl)
}

func TestServer_Create_EmptyURL(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: ""}

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Create_MissingHost(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: "http:///path"}

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Create_InvalidCharactersInHost(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: "https://exa mple.com"}

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Create_TooLongURL(t *testing.T) {
	server := &Server{}

	longURL := "https://example.com/" + strings.Repeat("a", 2000)
	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: longURL}

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Create_WhitespaceURL(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: " https://example.com "}

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Create_DuplicatedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: "https://example.com"}

	mockService.EXPECT().
		CreateURL(ctx, url.URL{OriginalURL: req.OriginalUrl}).
		Return(nil, url.ErrDuplicatedURL).
		Times(1)

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestServer_Create_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Create_Request{OriginalUrl: "https://example.com"}

	mockService.EXPECT().
		CreateURL(ctx, url.URL{OriginalURL: req.OriginalUrl}).
		Return(nil, errors.New("db error")).
		Times(1)

	resp, err := server.Create(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.Internal, status.Code(err))
}
