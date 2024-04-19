package service

import "scrapper/internal/domain/model"

type ScrapperService struct {
	scraper model.Scraper
}

func NewScrapperService(scraper model.Scraper) *ScrapperService {
	return &ScrapperService{
		scraper: scraper,
	}
}

func (s *ScrapperService) FetchData(url string) ([]byte, error) {
	return s.scraper.Scrape(url)
}
