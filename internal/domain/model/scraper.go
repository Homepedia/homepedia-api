package model

type Scraper interface {
	Scrape(url string) ([]byte, error)
}
