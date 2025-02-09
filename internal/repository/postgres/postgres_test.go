//go:build dbtest
// +build dbtest

package url

import (
	"context"
	"log"
	"os"
	"testing"

	r "github.com/kingxl111/url-shortener/internal/repository"
	"github.com/kingxl111/url-shortener/internal/url"
	"github.com/kingxl111/url-shortener/internal/url/shortener"
	"github.com/stretchr/testify/require"
)

var testDB *DB

func TestMain(m *testing.M) {
	db, err := NewDB("user", "password", "localhost", "5432", "shortener-db", "disable")
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	testDB = db
	defer testDB.Close()

	_, err = testDB.pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL UNIQUE,
			shortened_url TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		log.Fatalf("Cannot create table: %v", err)
	}

	code := m.Run()

	_, _ = testDB.pool.Exec(context.Background(), `TRUNCATE TABLE urls RESTART IDENTITY`)

	os.Exit(code)
}

func clearDB(t *testing.T) {
	_, err := testDB.pool.Exec(context.Background(), `TRUNCATE TABLE urls RESTART IDENTITY`)
	require.NoError(t, err, "Failed to clear database before test")
}

func TestRepository_Create(t *testing.T) {
	repo := NewRepository(testDB)
	ctx := context.Background()
	clearDB(t)

	original := "https://example.com"
	shortened := shortener.GenerateShortURL(original)
	newURL := url.URL{OriginalURL: original, ShortenedURL: shortened}

	createdURL, err := repo.Create(ctx, newURL)
	require.NoError(t, err, "Create() should not return an error")
	require.Equal(t, newURL.OriginalURL, createdURL.OriginalURL)
	require.Equal(t, newURL.ShortenedURL, createdURL.ShortenedURL)

	_, err = repo.Create(ctx, newURL)
	require.ErrorIs(t, err, r.ErrorDuplicatedURL, "Create() should return an error for duplicate URL")

	anotherURL := url.URL{
		OriginalURL:  "https://another-example.com",
		ShortenedURL: shortener.GenerateShortURL("https://another-example.com"),
	}
	createdAnother, err := repo.Create(ctx, anotherURL)
	require.NoError(t, err, "Create() should allow creating unique URLs")
	require.Equal(t, anotherURL.OriginalURL, createdAnother.OriginalURL)
}

func TestRepository_Get(t *testing.T) {
	repo := NewRepository(testDB)
	ctx := context.Background()
	clearDB(t)

	original := "https://example.com"
	shortened := shortener.GenerateShortURL(original)
	newURL := url.URL{OriginalURL: original, ShortenedURL: shortened}

	_, err := repo.Create(ctx, newURL)
	require.NoError(t, err, "Create() should not return an error")

	gotURL, err := repo.Get(ctx, url.URL{ShortenedURL: shortened})
	require.NoError(t, err, "Get() should return no error for existing URL")
	require.Equal(t, original, gotURL.OriginalURL)

	_, err = repo.Get(ctx, url.URL{ShortenedURL: "nonexistent123"})
	require.ErrorIs(t, err, r.ErrorNotFound, "Get() should return ErrNotFound for non-existent URL")

	_, err = repo.Get(ctx, url.URL{ShortenedURL: ""})
	require.Error(t, err, "Get() should return an error for empty short URL")
}

func TestRepository_Create_And_Get_With_Different_Formats(t *testing.T) {
	repo := NewRepository(testDB)
	ctx := context.Background()
	clearDB(t)

	testCases := []struct {
		original  string
		shortened string
	}{
		{"https://example.com", shortener.GenerateShortURL("https://example.com")},
		{"http://example.com", shortener.GenerateShortURL("http://example.com")},
		{"https://example.com/path?query=1", shortener.GenerateShortURL("https://example.com/path?query=1")},
		{"https://sub.example.com", shortener.GenerateShortURL("https://sub.example.com")},
		{"https://example.com/very/long/path/that/might/be/used", shortener.GenerateShortURL("https://example.com/very/long/path/that/might/be/used")},
	}

	for _, tc := range testCases {
		t.Run(tc.original, func(t *testing.T) {
			newURL := url.URL{OriginalURL: tc.original, ShortenedURL: tc.shortened}

			_, err := repo.Create(ctx, newURL)
			require.NoError(t, err, "Create() should not return error")

			gotURL, err := repo.Get(ctx, url.URL{ShortenedURL: tc.shortened})
			require.NoError(t, err, "Get() should not return error")
			require.Equal(t, tc.original, gotURL.OriginalURL)
		})
	}
}
