package scrapper

import (
	"homepedia-api/scraper/internal/cron"

	"go.mongodb.org/mongo-driver/mongo"
)

func InitService(client *mongo.Client) {
	cronService := cron.NewExecuteCron(client)
	cronService.Execute()
}
