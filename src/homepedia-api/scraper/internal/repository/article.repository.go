package repository

import (
	"context"
	"homepedia-api/scraper/internal/domain"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepository struct {
	Client *mongo.Client
}

type IArticleRepository interface {
	InsertMany(data *[]domain.FigaroData)
}

func NEewArticleRepository(client *mongo.Client) IArticleRepository {
	return &ArticleRepository{
		Client: client,
	}
}

func (e *ArticleRepository) InsertMany(data *[]domain.FigaroData) {
	collection := e.Client.Database("homepedia").Collection("figaroArticles")
	var models []mongo.WriteModel

	for _, d := range *data {
		filter := bson.M{"ref": d.Ref}
		update := bson.M{
			"$set": d,
		}
		upsert := true
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(upsert)
		models = append(models, model)
	}

	bulkOption := options.BulkWrite().SetOrdered(false)
	if len(models) > 0 {
		_, err := collection.BulkWrite(context.Background(), models, bulkOption)
		if err != nil {
			log.Printf("Error performing bulk write: %v\n", err)
		} else {
			log.Printf("Bulk write completed successfully\n")
		}
	} else {
		log.Printf("No data to write\n")
	}
}
