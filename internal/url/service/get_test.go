package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/repository/mocks"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetURL(t *testing.T) {
	repo := &mocks.URLRepository{}
	s := New(repo)

	validShortURL := "AbC123_789"
	validURL := &url.URL{
		ShortenedURL: validShortURL,
		OriginalURL:  "http://example.com",
	}

	tests := []struct {
		name           string
		input          url.URL
		mockReturn     *url.URL
		mockError      error
		expectedResult *url.URL
		expectedErr    error
	}{
		{
			name: "valid short URL",
			input: url.URL{
				ShortenedURL: validShortURL,
			},
			mockReturn:     validURL,
			expectedResult: validURL,
		},
		{
			name: "empty short URL",
			input: url.URL{
				ShortenedURL: "",
			},
			expectedErr: url.ErrEmptyURL,
		},
		{
			name: "invalid length1",
			input: url.URL{
				ShortenedURL: "abc",
			},
			expectedErr: url.ErrInvalidLength,
		},
		{
			name: "invalid length2",
			input: url.URL{
				ShortenedURL: "982nkjsd dj31ie ksjdf",
			},
			expectedErr: url.ErrInvalidLength,
		},
		{
			name: "invalid characters1",
			input: url.URL{
				ShortenedURL: "abcedf!234",
			},
			expectedErr: url.ErrInvalidCharacters,
		},
		{
			name: "invalid characters2",
			input: url.URL{
				ShortenedURL: "<bcedf111<",
			},
			expectedErr: url.ErrInvalidCharacters,
		},
		{
			name: "not found in repository",
			input: url.URL{
				ShortenedURL: validShortURL,
			},
			mockError:   repository.ErrorNotFound,
			expectedErr: repository.ErrorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.expectedErr == nil || errors.Is(tt.expectedErr, repository.ErrorNotFound) {
				repo.On("Get", mock.Anything, tt.input).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			result, err := s.GetURL(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			repo.AssertExpectations(t)
		})
	}
}

func TestIsValidShortURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid lowercase", "abcdefghij", true},
		{"valid uppercase", "ABCDEFGHIJ", true},
		{"valid with numbers", "a1b2c3d4e5", true},
		{"valid with underscore", "a_b_c_d_e_", true},
		{"invalid character", "abc@123", false},
		{"invalid space", "abc 123", false},
		{"invalid hyphen", "abc-123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidShortURL(tt.input))
		})
	}
}
