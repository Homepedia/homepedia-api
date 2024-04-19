package main

import (
	"fmt"
	"scrapper/internal/application/service"
	"scrapper/internal/infrastructure/external"
)

func main() {
	fmt.Println("Scrapper started...")

	scrapper := external.NewCollyScraper()
	collyService := service.NewScrapperService(scrapper)
	data, err := collyService.FetchData("https://www.google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

}
