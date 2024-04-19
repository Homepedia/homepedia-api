package service

import "scraper/internal/domain/model"

type ScraperService struct {
	scraper model.Scraper
}

func NewScrapperService(scraper model.Scraper) *ScraperService {
	return &ScraperService{
		scraper: scraper,
	}
}

func (s *ScraperService) FetchData(url string) ([]byte, error) {
	return s.scraper.Scrape(url)
}
