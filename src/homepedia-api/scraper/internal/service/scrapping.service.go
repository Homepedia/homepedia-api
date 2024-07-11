package service

import (
	"errors"
	"fmt"
	"homepedia-api/lib/utils"
	"homepedia-api/scraper/internal/domain"
	"homepedia-api/scraper/internal/repository"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/mongo"
)

func ScrappePage(url string, userAgentList []string, key string) (*domain.FigaroData, error) {
	log.Printf("Scrapping data from URL: %s\n", url)
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
		if city != "" {
			re := regexp.MustCompile(`à\s*|(\s*\(0[^)]*\))`)
			formattedCity := re.ReplaceAllString(city, "")
			data.City = formattedCity
		}
		data.Price = e.ChildText("div.classified-price-per-m2 > strong")
		energeticPerf := e.ChildText("ul > li.active.dpe-c")
		if len(energeticPerf) > 0 {
			data.EnergeticPerf = energeticPerf
		} else {
			data.EnergeticPerf = "N.C"
		}
		greenHouseGasIndex := e.ChildText("ul > li.active.ges-c")
		if len(greenHouseGasIndex) > 0 {
			data.GreenhouseGasIndex = greenHouseGasIndex
		} else {
			data.GreenhouseGasIndex = "N.C"
		}
		e.ForEach("ul.unstyled.features-list > li", func(_ int, el *colly.HTMLElement) {
			class := el.ChildAttr("i", "class")
			switch class {
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
		keys := strings.Split(key, "-")
		if len(keys) >= 2 {
			data.Region = keys[0]
			data.Department = keys[1]
		}
		data.Source = "Figaro Immobilier"
		data.Type = 1

		log.Printf("Scrapped data: %+v\n", data)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error on URL %s: %v\n", url, err)
	})

	log.Printf("Visiting URL: %s\n", url)
	c.Visit(url)
	c.Wait()

	log.Printf("Finished scrapping URL: %s\n", url)

	return &data, nil
}

func ScrapeDepartment(departmentURL, key string, userAgentList []string, maxCount int, wg *sync.WaitGroup, mu *sync.Mutex, data *[]domain.FigaroData, sem chan struct{}, client *mongo.Client) {
	defer wg.Done()
	var links []string
	pageCounter := 1
	visitedURLs := make(map[string]bool)
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("immobilier.lefigaro.fr"),
	)

	c.SetRequestTimeout(120 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*immobilier.lefigaro.fr*",
		Delay:       2 * time.Second,
		Parallelism: 3,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.RandomString(userAgentList))
		time.Sleep(5 * time.Second)
		log.Printf("Visiting URL: %s\n", r.URL.String())
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
			log.Printf("Visiting next page: %s\n", nextPage)
			c.Visit(nextPage)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v\n", err)
		if r.StatusCode == 429 {
			log.Println("Received 429, waiting before retrying...")
			time.Sleep(30 * time.Second)
			r.Request.Retry()
		}
	})

	log.Printf("Visiting department URL: %s\n", departmentURL)
	c.Visit(departmentURL)
	c.Wait()
	for _, link := range links {
		res, err := ScrappePage(link, userAgentList, key)
		if err == nil {
			mu.Lock()
			*data = append(*data, *res)
			mu.Unlock()
		} else {
			log.Printf("Error scraping URL %s: %v\n", link, err)
		}
	}

	log.Printf("Finished scrapping department: %s\n", departmentURL)
	repository := repository.NEewArticleRepository(client)
	repository.InsertMany(data)

	<-sem
}
