package url

import "errors"

var (
	ErrEmptyURL      = errors.New("URL is empty")
	ErrInvalidFormat = errors.New("invalid URL format")
	ErrMissingHost   = errors.New("missing host")
	ErrRepository    = errors.New("repository error")
	ErrDuplicatedURL = errors.New("duplicated URL")
	ErrService       = errors.New("service error")
	ErrNotFoundURL   = errors.New("url not found")
)
