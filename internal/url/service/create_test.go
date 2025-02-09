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

func TestService_CreateURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	inputURL := url.URL{OriginalURL: "https://example.com"}
	expectedShortened := "abc123"

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(&url.URL{OriginalURL: inputURL.OriginalURL, ShortenedURL: expectedShortened}, nil).
		Times(1)

	result, err := svc.CreateURL(ctx, inputURL)
	require.NoError(t, err)
	require.Equal(t, expectedShortened, result.ShortenedURL)
	require.Equal(t, inputURL.OriginalURL, result.OriginalURL)
}

func TestService_CreateURL_Duplicated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	inputURL := url.URL{OriginalURL: "https://example.com"}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil, repository.ErrorDuplicatedURL).
		Times(1)

	result, err := svc.CreateURL(ctx, inputURL)
	require.ErrorIs(t, err, url.ErrDuplicatedURL)
	require.Nil(t, result)
}

func TestService_CreateURL_UnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLRepository(ctrl)
	svc := New(mockRepo)

	ctx := context.Background()
	inputURL := url.URL{OriginalURL: "https://example.com"}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil, errors.New("db failure")).
		Times(1)

	result, err := svc.CreateURL(ctx, inputURL)
	require.ErrorIs(t, err, url.ErrService)
	require.Nil(t, result)
}
