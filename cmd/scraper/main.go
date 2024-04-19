package main

import (
	"fmt"
	"scraper/internal/application/service"
	"scraper/internal/infrastructure/external"
)

func main() {
	fmt.Println("Scraper started...")

	scraper := external.NewCollyScraper()
	collyService := service.NewScrapperService(scraper)
	data, err := collyService.FetchData("https://www.google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

}
