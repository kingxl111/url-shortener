package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kingxl111/url-shortener/internal/url"
	shrt "github.com/kingxl111/url-shortener/pkg/shortener"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: "abc1234567"}
	expectedOriginalURL := "https://example.com"

	mockService.EXPECT().
		GetURL(ctx, url.URL{ShortenedURL: req.ShortUrl}).
		Return(&url.URL{OriginalURL: expectedOriginalURL, ShortenedURL: req.ShortUrl}, nil).
		Times(1)

	resp, err := server.Get(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expectedOriginalURL, resp.OriginalUrl)
}

func TestServer_Get_EmptyShortURL(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: ""}

	resp, err := server.Get(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Get_InvalidShortURLLength(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: "abc"}

	resp, err := server.Get(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Get_InvalidShortURLFormat(t *testing.T) {
	server := &Server{}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: "abc@123123"}

	resp, err := server.Get(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: "1231abc123"}

	mockService.EXPECT().
		GetURL(ctx, url.URL{ShortenedURL: req.ShortUrl}).
		Return(nil, url.ErrNotFoundURL).
		Times(1)

	resp, err := server.Get(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestServer_Get_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockShortenerService(ctrl)
	server := &Server{Services: mockService}

	ctx := context.Background()
	req := &shrt.Get_Request{ShortUrl: "a1231bc123"}

	mockService.EXPECT().
		GetURL(ctx, url.URL{ShortenedURL: req.ShortUrl}).
		Return(nil, errors.New("db error")).
		Times(1)

	resp, err := server.Get(ctx, req)
	require.Nil(t, resp)
	require.Equal(t, codes.Internal, status.Code(err))
}
