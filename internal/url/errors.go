package url

import "errors"

var (
	ErrEmptyURL      = errors.New("URL is empty")
	ErrInvalidFormat = errors.New("invalid URL format")
	ErrMissingHost   = errors.New("missing host")
)
