package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/stretchr/testify/require"
)

func TestService_GetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	shortURL := url.URL{ShortenedURL: "abc1234567"}
	expectedOriginal := "https://example.com"

	mockRepo.EXPECT().
		Get(ctx, shortURL).
		Return(&url.URL{OriginalURL: expectedOriginal, ShortenedURL: shortURL.ShortenedURL}, nil).
		Times(1)

	result, err := svc.GetURL(ctx, shortURL)
	require.NoError(t, err)
	require.Equal(t, expectedOriginal, result.OriginalURL)
	require.Equal(t, shortURL.ShortenedURL, result.ShortenedURL)
}

func TestService_GetURL_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	shortURL := url.URL{ShortenedURL: "4567abc123"}

	mockRepo.EXPECT().
		Get(ctx, shortURL).
		Return(nil, repository.ErrorNotFound).
		Times(1)

	result, err := svc.GetURL(ctx, shortURL)
	require.ErrorIs(t, err, url.ErrNotFoundURL)
	require.Nil(t, result)
}

func TestService_GetURL_UnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	shortURL := url.URL{ShortenedURL: "4567abc123"}

	mockRepo.EXPECT().
		Get(ctx, shortURL).
		Return(nil, errors.New("db failure")).
		Times(1)

	result, err := svc.GetURL(ctx, shortURL)
	require.ErrorIs(t, err, url.ErrRepository)
	require.Nil(t, result)
}
