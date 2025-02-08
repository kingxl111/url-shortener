package service

import (
	"context"
	"errors"
	"github.com/kingxl111/url-shortener/internal/repository/mocks"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
	"strings"
	"testing"
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
	repo := &mocks.URLRepository{}
	s := New(repo)

	tests := []struct {
		name           string
		input          url.URL
		mockRepoError  error
		expectedErr    error
		expectedLength int
	}{
		{
			name: "successful creation",
			input: url.URL{
				OriginalURL: "example.com",
			},
			expectedLength: shortener.ShortURLLength,
		},
		{
			name: "invalid url format",
			input: url.URL{
				OriginalURL: "http://invalid!host",
			},
			expectedErr: url.ErrInvalidFormat,
		},
		{
			name: "repository error (e.g., duplicate)",
			input: url.URL{
				OriginalURL: "example.com",
			},
			mockRepoError: errors.New("duplicate key"),
			expectedErr:   errors.New("repository error: duplicate key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настройка мока
			repo.On("Create", context.Background(), tt.input).
				Return(&tt.input, tt.mockRepoError).
				Once()

			result, err := s.CreateURL(context.Background(), tt.input)

			// Проверка ошибок
			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Проверка длины короткого URL
			if len(result.ShortenedURL) != tt.expectedLength {
				t.Fatalf("expected length: %d, got: %d", tt.expectedLength, len(result.ShortenedURL))
			}

			// Проверка символов короткого URL
			for _, c := range result.ShortenedURL {
				if !isValidShortURLChar(c) {
					t.Fatalf("invalid character in short URL: %c", c)
				}
			}

			// Проверка нормализации оригинального URL
			if result.OriginalURL != "http://example.com" {
				t.Fatalf("normalization failed, got: %s", result.OriginalURL)
			}
		})
	}
}

func isValidShortURLChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
