package shortener

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateShortURL(t *testing.T) {
	tests := []struct {
		name      string
		original  string
		wantLen   int
		wantEqual bool
	}{
		{"Normal URL", "https://example.com", ShortURLLength, true},
		{"URL with path", "https://example.com/path", ShortURLLength, true},
		{"URL with query", "https://example.com?q=123", ShortURLLength, true},
		{"Different URL", "https://another.com", ShortURLLength, false},
		{"Empty URL", "", ShortURLLength, true},
	}

	shortURLs := make(map[string]string) // Для проверки уникальности

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURL := GenerateShortURL(tt.original)

			require.Equal(t, tt.wantLen, len(shortURL), "Shortened URL length should be %d", tt.wantLen)

			require.True(t, IsValidShortURL(shortURL), "Shortened URL contains invalid characters: %s", shortURL)

			if prev, exists := shortURLs[tt.original]; exists {
				require.Equal(t, prev, shortURL, "Same input should produce the same short URL")
			} else {
				shortURLs[tt.original] = shortURL
			}
		})
	}
}

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name string
		num  uint64
		want int
	}{
		{"Encode 0", 0, ShortURLLength},
		{"Encode 1", 1, ShortURLLength},
		{"Encode max uint64", ^uint64(0), ShortURLLength},
		{"Encode small number", 42, ShortURLLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeBase62(tt.num)
			require.Equal(t, tt.want, len(got), "Encoded string should have fixed length")
			require.True(t, IsValidShortURL(got), "Encoded string contains invalid characters")
		})
	}
}
