package cron

import (
	"fmt"

	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExecuteCron struct {
	Client *mongo.Client
}

type IExecuteCron interface {
	Execute()
}

func NewExecuteCron(client *mongo.Client) IExecuteCron {
	return &ExecuteCron{
		Client: client,
	}
}
func (e *ExecuteCron) Execute() {
	c := cron.New()
	c.AddFunc("0 02 17 * * ?", func() {
		fmt.Println("Scraping started..")
		if e.Client != nil {
			fmt.Println("Client MongoDB is ready.")
			RunFigaroCron(e.Client)
		} else {
			fmt.Println("Client MongoDB is nil. Cannot start scraping.")
		}
	})

	fmt.Println("Cron scheduler started. Waiting for scraping job to start...")

	c.Start()

	select {}
}
