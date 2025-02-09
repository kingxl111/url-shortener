package url

import "errors"

var (
	ErrEmptyURL      = errors.New("URL is empty")
	ErrInvalidFormat = errors.New("invalid URL format")
	ErrMissingHost   = errors.New("missing host")

	ErrInvalidCharacters = errors.New("invalid characters")
	ErrInvalidLength     = errors.New("invalid length")
	ErrRepository        = errors.New("repository error")
	ErrInvalidScheme     = errors.New("invalid scheme")
)
