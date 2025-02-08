package in_memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/kingxl111/url-shortener/internal/model"
	repos "github.com/kingxl111/url-shortener/internal/repository"
)

var _ repos.URLRepository = (*MemoryStorage)(nil)

type MemoryStorage struct {
	mu              sync.RWMutex
	shortToOriginal map[string]string
	originalToShort map[string]string
}

func NewMemoryStorage() repos.URLRepository {
	return &MemoryStorage{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

func (m *MemoryStorage) Create(ctx context.Context, url model.URL) (model.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.originalToShort[url.OriginalURL]; exist {
		return url, fmt.Errorf("url already exists: %s", url.OriginalURL)
	}
	m.shortToOriginal[url.ShortenedURL] = url.OriginalURL
	m.originalToShort[url.OriginalURL] = url.ShortenedURL

	return url, nil
}

func (m *MemoryStorage) Get(ctx context.Context, url model.URL) (model.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	original, exists := m.shortToOriginal[url.ShortenedURL]
	if !exists {
		return url, fmt.Errorf("such url does not exist: %s", url.ShortenedURL)
	}
	url.OriginalURL = original

	return url, nil
}
