package in_memory

import (
	"context"
	"sync"

	"github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/url"
)

type MemoryStorage struct {
	mu              sync.RWMutex
	shortToOriginal map[string]string
	originalToShort map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		shortToOriginal: make(map[string]string),
		originalToShort: make(map[string]string),
	}
}

func (m *MemoryStorage) Create(ctx context.Context, url url.URL) (*url.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if _, exist := m.originalToShort[url.OriginalURL]; exist {
			return nil, repository.ErrorDuplicatedURL
		}
		m.shortToOriginal[url.ShortenedURL] = url.OriginalURL
		m.originalToShort[url.OriginalURL] = url.ShortenedURL

		return &url, nil
	}
}

func (m *MemoryStorage) Get(ctx context.Context, url url.URL) (*url.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		original, exists := m.shortToOriginal[url.ShortenedURL]
		if !exists {
			return nil, repository.ErrorNotFound
		}
		url.OriginalURL = original

		return &url, nil
	}

}
