package main

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

// RandomString retourne une chaîne aléatoire de la liste fournie
func RandomString(userAgentList []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAgentList))
	return userAgentList[randomIndex]
}

// FigaroData structure to hold scraped data
type FigaroData struct {
	Price              string `json:"price"`
	City               string `json:"city"`
	EnergeticPerf      string `json:"dpe"`
	GreenhouseGasIndex string `json:"ges"`
	Bedrooms           string `json:"bedrooms"`
	Bathrooms          string `json:"bathrooms"`
	Rooms              string `json:"rooms"`
	Area               string `json:"area"`
	Url                string `json:"url"`
	Source             string `json:"source"`
	Ref                string `json:"ref"`
	Region             string `json:"region"`
	Department         string `json:"department"`
	Type               int    `json:"type"`
}

// ScrappePage scrapes the data from a given URL
func ScrappePage(url string, userAgentList []string) (*FigaroData, error) {
	re := regexp.MustCompile(`annonce-(\d+)`)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return nil, errors.New("cannot get id from url")
	}
	id := match[1]
	var data FigaroData

	data.Ref = id

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("immobilier.lefigaro.fr"),
		colly.UserAgent(RandomString(userAgentList)),
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
		data.Source = "Figaro Immobilier"
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erreur sur le lien %s: %v\n", url, err)
	})

	c.Visit(url)
	c.Wait()

	return &data, nil
}

// ScrapeDepartment scrapes all pages for a given department
func ScrapeDepartment(departmentURL, departmentCode string, userAgentList []string, maxCount int, wg *sync.WaitGroup, mu *sync.Mutex, data *[]FigaroData, sem chan struct{}) {
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
		r.Headers.Set("User-Agent", RandomString(userAgentList))
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
		res, err := ScrappePage(link, userAgentList)
		if err == nil {
			mu.Lock()
			*data = append(*data, *res)
			mu.Unlock()
		}
	}

	<-sem // Libère une place dans le canal
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data []FigaroData
	sem := make(chan struct{}, 5)

	departments := map[string]string{
		"55-04": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-alpes+de+haute+provence.html",
		"55-05": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-hautes+alpes.html",
		"55-06": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-alpes+maritimes.html",
		"55-13": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-bouches+du+rhone.html",
		"55-83": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-var.html",
		"55-84": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-vaucluse.html",
		"44-67": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-bas+rhin.html",
		"44-68": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haut+rhin.html",
		"75-24": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-dordogne.html",
		"75-33": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-gironde.html",
		"75-40": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-landes.html",
		"75-47": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-lot+et+garonne.html",
		"75-64": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-pyrenees+atlantiques.html",
		"84-03": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-allier.html",
		"84-15": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-cantal.html",
		"84-43": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+loire.html",
		"84-63": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-puy+de+dome.html",
		"28-14": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-calvados.html",
		"28-50": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-manche.html",
		"28-61": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-orne.html",
		"27-21": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-cote+d+or.html",
		"27-58": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-nievre.html",
		"27-71": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-saone+et+loire.html",
		"27-89": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-yonne.html",
		"53-22": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-cotes+d+armor.html",
		"53-29": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-finistere.html",
		"53-35": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-ille+et+vilaine.html",
		"53-56": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-morbihan.html",
		"24-18": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-cher.html",
		"24-28": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-eure+et+loir.html",
		"24-36": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-indre.html",
		"24-37": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-indre+et+loire.html",
		"24-41": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-loir+et+cher.html",
		"21-45": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-loiret.html",
		"44-08": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-ardennes.html",
		"44-10": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-aube.html",
		"44-51": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-marne.html",
		"44-52": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+marne.html",
		"94-2A": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-corse+du+sud.html",
		"94-2B": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+corse.html",
		"27-25": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-doubs.html",
		"27-39": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-jura.html",
		"27-70": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+saone.html",
		"27-90": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-territoire+de+belfort.html",
		"28-27": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-eure.html",
		"28-76": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-seine+maritime.html",
		"11-75": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-paris.html",
		"11-77": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-seine+et+marne.html",
		"11-78": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-yvelines.html",
		"11-91": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-essonne.html",
		"11-92": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-hauts+de+seine.html",
		"11-93": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-seine+saint+denis.html",
		"11-94": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-val+de+marne.html",
		"11-95": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-val+d+oise.html",
		"76-11": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-aude.html",
		"76-30": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-gard.html",
		"76-34": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-herault.html",
		"76-48": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-lozere.html",
		"76-66": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-pyrenees+orientales.html",
		"75-19": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-correze.html",
		"75-23": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-creuse.html",
		"75-87": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+vienne.html",
		"44-54": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-meurthe+et+moselle.html",
		"44-55": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-meuse.html",
		"44-57": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-moselle.html",
		"44-88": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-vosges.html",
		"76-09": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-ariege.html",
		"76-12": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-aveyron.html",
		"76-31": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+garonne.html",
		"76-32": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-gers.html",
		"76-46": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-lot.html",
		"76-65": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-hautes+pyr%C3%A9n%C3%A9es.html",
		"76-81": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-tarn.html",
		"76-82": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-tarn+et+garonne.html",
		"32-59": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-nord.html",
		"32-62": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-pas+de+calais.html",
		"52-44": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-loire+atlantique.html",
		"52-49": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-maine+et+loire.html",
		"52-53": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-mayenne.html",
		"52-72": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-sarthe.html",
		"52-85": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-vendee.html",
		"32-02": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-aisne.html",
		"32-60": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-oise.html",
		"32-80": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-somme.html",
		"75-16": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-charente.html",
		"75-17": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-charente+maritime.html",
		"75-79": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-deux+sevres.html",
		"75-86": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-vienne.html",
		"84-01": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-ain.html",
		"84-07": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-ardeche.html",
		"84-26": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-drome.html",
		"84-38": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-isere.html",
		"84-42": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-loire.html",
		"84-69": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-rhone+alpes.html",
		"84-73": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-savoie.html",
		"84-74": "https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-haute+savoie.html",
	}
	// departments := []string{
	// 	"https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-alpes+maritimes.html",
	// 	"https://immobilier.lefigaro.fr/annonces/immobilier-vente-bien-rhone.html",
	// }

	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
	}

	maxCount := 100

	for departement, departmentURL := range departments {
		sem <- struct{}{} // Bloque si le canal est plein
		wg.Add(1)
		go ScrapeDepartment(departmentURL, departement, userAgentList, maxCount, &wg, &mu, &data, sem)
	}

	wg.Wait()

	// Afficher les données collectées
	for _, d := range data {
		fmt.Printf("%+v\n", d)
	}
}
