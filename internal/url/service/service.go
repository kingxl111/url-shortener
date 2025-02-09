package service

type service struct {
	urlRepository URLRepository
}

func New(urlRepository URLRepository) *service {
	return &service{
		urlRepository: urlRepository,
	}
}
