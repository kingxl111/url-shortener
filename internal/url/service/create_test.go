package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/repository/mocks"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
	"github.com/stretchr/testify/mock"
)

func TestValidateAndNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:     "valid http url with path",
			input:    "http://example.com/path",
			expected: "http://example.com/path",
		},
		{
			name:     "valid https url with query",
			input:    "https://example.com?param=value",
			expected: "https://example.com?param=value",
		},
		{
			name:     "add http scheme if missing",
			input:    "example.com",
			expected: "http://example.com",
		},
		{
			name:        "empty url",
			input:       "",
			expectedErr: url.ErrEmptyURL,
		},
		{
			name:        "invalid url format",
			input:       "http://%invalid.com",
			expectedErr: url.ErrInvalidFormat,
		},
		{
			name:        "missing host",
			input:       "http://",
			expectedErr: url.ErrMissingHost,
		},
		{
			name:     "normalize scheme and host to lowercase",
			input:    "HTTPS://EXAMPLE.COM/Path",
			expected: "https://example.com/Path",
		},

		{
			name:     "url with port",
			input:    "http://example.com:8080",
			expected: "http://example.com:8080",
		},
		{
			name:     "url with port and path",
			input:    "example.com:8080/path",
			expected: "http://example.com:8080/path",
		},

		{
			name:        "ftp scheme (if allowed)",
			input:       "ftp://example.com/file.txt",
			expectedErr: url.ErrInvalidScheme,
		},
		{
			name:     "url with query and fragment",
			input:    "http://example.com/path?param=val#anchor",
			expected: "http://example.com/path?param=val#anchor",
		},

		{
			name:     "url with spaces",
			input:    "  http://example.com  ",
			expected: "http://example.com",
		},
		{
			name:        "url with invalid characters in host",
			input:       "http://exa mple.com",
			expectedErr: url.ErrInvalidFormat,
		},
		{
			name:     "url with encoded spaces",
			input:    "http://example.com/%20path",
			expected: "http://example.com/%20path",
		},

		{
			name:        "websocket scheme",
			input:       "ws://example.com/ws",
			expectedErr: url.ErrInvalidScheme,
		},

		{
			name:     "subdomain",
			input:    "http://api.example.com",
			expected: "http://api.example.com",
		},
		{
			name:     "multiple subdomains",
			input:    "http://v1.api.example.com",
			expected: "http://v1.api.example.com",
		},

		{
			name:     "punycode domain",
			input:    "http://xn--e1afmkfd.xn--p1ai",
			expected: "http://xn--e1afmkfd.xn--p1ai",
		},

		{
			name:     "long url",
			input:    "http://example.com/" + strings.Repeat("a", 1000),
			expected: "http://example.com/" + strings.Repeat("a", 1000),
		},
		{
			name:     "minimal valid url",
			input:    "http://a.b",
			expected: "http://a.b",
		},
		{
			name:        "only scheme",
			input:       "http://",
			expectedErr: url.ErrMissingHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndNormalizeURL(tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if result != tt.expected {
				t.Fatalf("expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	repo := new(mocks.URLRepository)
	s := New(repo)

	tests := []struct {
		name           string
		input          url.URL
		expectedInput  url.URL
		mockRepoError  error
		expectedErr    error
		expectedLength int
	}{
		{
			name: "successful creation",
			input: url.URL{
				OriginalURL: "http://example.com",
			},
			expectedInput: url.URL{
				OriginalURL:  "http://example.com",
				ShortenedURL: shortener.GenerateShortURL("http://example.com"),
			},
			expectedLength: shortener.ShortURLLength,
		},
		{
			name: "invalid url format",
			input: url.URL{
				OriginalURL: "http://invalid!host",
			},
			expectedErr:    url.ErrInvalidFormat,
			expectedLength: 0,
		},
		{
			name: "duplicate url",
			input: url.URL{
				OriginalURL: "http://duplicate.com",
			},
			expectedInput: url.URL{
				OriginalURL:  "http://duplicate.com",
				ShortenedURL: shortener.GenerateShortURL("http://duplicate.com"),
			},
			mockRepoError:  repository.ErrorDuplicatedURL,
			expectedErr:    errors.New("url already exists"),
			expectedLength: 0,
		},
		{
			name: "repository failure",
			input: url.URL{
				OriginalURL: "http://example.org",
			},
			mockRepoError:  errors.New("database error"),
			expectedErr:    errors.New("internal server error"),
			expectedLength: 0,
		},
		{
			name: "empty url",
			input: url.URL{
				OriginalURL: "",
			},
			expectedErr:    url.ErrEmptyURL,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.expectedErr == nil {
				repo.On("Create", ctx, mock.MatchedBy(func(u url.URL) bool {
					return u.OriginalURL == tt.expectedInput.OriginalURL
				})).Return(&tt.expectedInput, nil).Once()
			} else if tt.mockRepoError != nil {
				repo.On("Create", ctx, mock.MatchedBy(func(u url.URL) bool {
					return u.OriginalURL == tt.input.OriginalURL
				})).Return(nil, tt.mockRepoError).Once()
			}

			result, err := s.CreateURL(ctx, tt.input)

			if tt.expectedErr != nil {
				if err == nil || tt.expectedErr.Error() != err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
				}
				repo.AssertNotCalled(t, "Create")
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.ShortenedURL) != tt.expectedLength {
				t.Fatalf("expected length: %d, got: %d", tt.expectedLength, len(result.ShortenedURL))
			}

			if result.OriginalURL != tt.expectedInput.OriginalURL {
				t.Fatalf("expected original URL: %s, got: %s", tt.expectedInput.OriginalURL, result.OriginalURL)
			}

			repo.AssertExpectations(t)
		})
	}
}
