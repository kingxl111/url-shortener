package service

import (
	"github.com/kingxl111/url-shortener/internal/repository"
)

type service struct {
	urlRepository repository.URLRepository
}

func New(urlRepository repository.URLRepository) *service {
	return &service{
		urlRepository: urlRepository,
	}
}
