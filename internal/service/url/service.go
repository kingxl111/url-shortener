package url

import (
	"github.com/kingxl111/url-shortener/internal/repository"
	srv "github.com/kingxl111/url-shortener/internal/service"
)

var _ srv.ShortenerService = (*service)(nil)

type service struct {
	urlRepository repository.URLRepository
}

func New(urlRepository repository.URLRepository) srv.ShortenerService {
	return &service{
		urlRepository: urlRepository,
	}
}
