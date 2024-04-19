package external

import (
	"github.com/gocolly/colly/v2"
)

type CollyScraper struct {
	collector *colly.Collector
}

func NewCollyScraper() *CollyScraper {
	return &CollyScraper{
		collector: colly.NewCollector(),
	}
}

func (c *CollyScraper) Scrape(url string) ([]byte, error) {
	var result []byte
	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		result = e.Response.Body
	})

	err := c.collector.Visit(url)
	if err != nil {
		return nil, err
	}

	return result, nil
}
