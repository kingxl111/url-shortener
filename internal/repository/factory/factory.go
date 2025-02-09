package factory

import (
	"fmt"
	"github.com/kingxl111/url-shortener/internal/url/service"
	"os"

	m "github.com/kingxl111/url-shortener/internal/repository/in-memory"
	p "github.com/kingxl111/url-shortener/internal/repository/postgres"
)

const (
	memory   = "memory"
	database = "postgres"
)

func NewURLRepository(username, password, host, port, dbName, sslMode string) (service.URLRepository, error) {
	storageType := os.Getenv("STORAGE_TYPE")

	switch storageType {
	case memory:
		return m.NewMemoryStorage(), nil
	case database:
		db, err := p.NewDB(
			username,
			password,
			host,
			port,
			dbName,
			sslMode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init DB: %w", err)
		}
		return p.NewRepository(db), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
