package factory

import (
	"fmt"
	"os"

	"github.com/kingxl111/url-shortener/internal/repository"
	inmemory "github.com/kingxl111/url-shortener/internal/repository/in-memory"
	postgres "github.com/kingxl111/url-shortener/internal/repository/postgres"
)

func NewURLRepository(username, password, host, port, dbName, sslMode string) (repository.URLRepository, error) {
	storageType := os.Getenv("STORAGE_TYPE")

	switch storageType {
	case "memory":
		return inmemory.NewMemoryStorage(), nil
	case "postgres":
		db, err := postgres.NewDB(
			username,
			password,
			host,
			port,
			dbName,
			sslMode)

		if err != nil {
			return nil, fmt.Errorf("failed to init DB: %w", err)
		}
		return postgres.NewRepository(db), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
