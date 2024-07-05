package cron

import "go.mongodb.org/mongo-driver/mongo"

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
	RunFigaroCron(e.Client)
}
