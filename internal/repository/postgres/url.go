package url

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	rep "github.com/kingxl111/url-shortener/internal/repository"
	ur "github.com/kingxl111/url-shortener/internal/url"
)

const (
	tableName = "urls"

	urlColumn          = "original_url"
	shortenedURLColumn = "shortened_url"
)

type repository struct {
	db *DB
}

func NewRepository(db *DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, url ur.URL) (*ur.URL, error) {

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(urlColumn, shortenedURLColumn).
		Values(url.OriginalURL, url.ShortenedURL).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building insert query error: %w", err)
	}

	var id int64
	err = r.db.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, rep.ErrorDuplicatedURL
	}

	return &url, nil
}

func (r *repository) Get(ctx context.Context, shortenedUrl ur.URL) (*ur.URL, error) {

	builder := sq.Select(urlColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{shortenedURLColumn: shortenedUrl.ShortenedURL})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query error: %w", err)
	}

	row := r.db.pool.QueryRow(ctx, query, args...)

	err = row.Scan(&shortenedUrl.OriginalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, rep.ErrorNotFound
	} else if err != nil {
		return nil, fmt.Errorf("executing select query error: %w", err)
	}

	return &shortenedUrl, nil
}
