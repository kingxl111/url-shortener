package url

import (
	"context"
	"fmt"

	"github.com/kingxl111/url-shortener/internal/model"
	repo "github.com/kingxl111/url-shortener/internal/repository"

	sq "github.com/Masterminds/squirrel"
)

var _ repo.URLRepository = (*repository)(nil)

const (
	tableName = "urls"

	idColumn           = "id"
	urlColumn          = "original_url"
	shortenedURLColumn = "shortened_url"
)

type repository struct {
	db *DB
}

func NewRepository(db *DB) repo.URLRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, url model.URL) (model.URL, error) {

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(urlColumn, shortenedURLColumn).
		Values(url.OriginalURL, url.ShortenedURL).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return model.URL{}, fmt.Errorf("building insert query error: %w", err)
	}

	var id int64
	err = r.db.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return model.URL{}, fmt.Errorf("executing query error: %w", err)
	}
	url.ID = id

	return url, nil
}

func (r *repository) Get(ctx context.Context, shortenedUrl model.URL) (model.URL, error) {

	builder := sq.Select(urlColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{shortenedURLColumn: shortenedUrl.ShortenedURL})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.URL{}, fmt.Errorf("building select query error: %w", err)
	}

	row := r.db.pool.QueryRow(ctx, query, args...)
	err = row.Scan(&shortenedUrl.OriginalURL)
	if err != nil {
		return model.URL{}, fmt.Errorf("executing select query error: %w", err)
	}

	return shortenedUrl, nil
}
