package repository

import "errors"

var (
	ErrorNotFound      = errors.New("url not found")
	ErrorDuplicatedURL = errors.New("duplicated url")
)
