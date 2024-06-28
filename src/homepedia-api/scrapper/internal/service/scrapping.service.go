package service

import (
	"errors"
	"fmt"
	"homepedia-api/lib/utils"
	"homepedia-api/scrapper/internal/domain"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

// ScrappePage scrapes the data from a given URL
func ScrappePage(url string, userAgentList []string, key string) (*domain.FigaroData, error) {
	re := regexp.MustCompile(`annonce-(\d+)`)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return nil, errors.New("cannot get id from url")
	}
	id := match[1]
	var data domain.FigaroData

	data.Ref = id

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("immobilier.lefigaro.fr"),
		colly.UserAgent(utils.RandomString(userAgentList)),
		colly.MaxDepth(1),
	)

	c.OnHTML("div.classified-content", func(e *colly.HTMLElement) {
		city := e.ChildText("div.classified-main-infos-title > h1 > span")
		re := regexp.MustCompile(`à\s*|(\s*\(0[^)]*\))`)
		formattedCity := re.ReplaceAllString(city, "")
		data.City = formattedCity
		data.EnergeticPerf = e.ChildText("ul > li.active.dpe-c")
		data.GreenhouseGasIndex = e.ChildText("ul > li.active.ges-c")

		e.ForEach("ul.unstyled.features-list > li", func(_ int, el *colly.HTMLElement) {
			switch el.ChildAttr("i", "class") {
			case "spr-detail ic-area":
				data.Area = el.ChildText("span.feature")
			case "spr-detail ic-room":
				data.Rooms = el.ChildText("span.feature")
			case "spr-detail ic-bedroom":
				data.Bedrooms = el.ChildText("span.feature")
			case "spr-detail ic-bathroom":
				data.Bathrooms = el.ChildText("span.feature")
			}
		})

		data.Url = url
		data.Region = strings.Split(key, ",")[0]
		data.Department = strings.Split(key, ",")[1]
		data.Source = "Figaro Immobilier"
		data.Type = 1
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erreur sur le lien %s: %v\n", url, err)
	})

	c.Visit(url)
	c.Wait()

	return &data, nil
}

// ScrapeDepartment scrapes all pages for a given department
func ScrapeDepartment(departmentURL, key string, userAgentList []string, maxCount int, wg *sync.WaitGroup, mu *sync.Mutex, data *[]domain.FigaroData, sem chan struct{}) {
	defer wg.Done()
	var links []string
	pageCounter := 1
	visitedURLs := make(map[string]bool)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("immobilier.lefigaro.fr"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomString(userAgentList))
	})

	c.OnHTML("ul.list-annonce", func(e *colly.HTMLElement) {
		e.ForEach("a.content__link", func(_ int, el *colly.HTMLElement) {
			link := el.Attr("href")
			if strings.HasPrefix(link, "/annonces") {
				fullLink := "https://immobilier.lefigaro.fr" + link
				if !visitedURLs[fullLink] {
					visitedURLs[fullLink] = true
					links = append(links, fullLink)
				}
			}
		})
	})

	c.OnHTML("a.pagination__link", func(e *colly.HTMLElement) {
		if pageCounter < maxCount {
			pageCounter++
			nextPage := fmt.Sprintf("%s?page=%d", departmentURL, pageCounter)
			c.Visit(nextPage)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erreur: %v\n", err)
	})

	c.Visit(departmentURL)
	c.Wait()

	for _, link := range links {
		res, err := ScrappePage(link, userAgentList, key)
		if err == nil {
			mu.Lock()
			*data = append(*data, *res)
			mu.Unlock()
		}
	}

	<-sem // Libère une place dans le canal
}
