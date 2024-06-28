package cron

import "go.mongodb.org/mongo-driver/mongo"

type ExecuteCron struct {
	Client *mongo.Client
}

type IExecuteCron interface {
	Execute()
}

func (e *ExecuteCron) Execute() {}
