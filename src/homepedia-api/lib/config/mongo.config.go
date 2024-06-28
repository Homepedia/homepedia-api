package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDBConfig initialise la connexion à la base de données MongoDB
func InitMongoDBConfig() (*mongo.Client, error) {
	// Récupération des variables d'environnement
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	// Vérification des variables d'environnement
	if user == "" || password == "" || host == "" || port == "" {
		return nil, fmt.Errorf("les variables d'environnement MONGO_USER, MONGO_PASSWORD, MONGO_HOST et MONGO_PORT doivent être définies")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du client MongoDB: %v", err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, fmt.Errorf("erreur lors du ping à MongoDB: %v", err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	databaseNames, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la liste des bases de données: %v", err)
	}

	dbExists := false
	for _, name := range databaseNames {
		if name == "homepedia" {
			dbExists = true
			break
		}
	}

	if !dbExists {
		err = client.Database("homepedia").CreateCollection(context.TODO(), "figaroArticles")
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la création de la collection 'figaroArticles' dans la nouvelle base de données 'homepedia': %v", err)
		}
		fmt.Println("Base de données 'homepedia' et collection 'figaroArticles' créées avec succès.")
	} else {
		collectionNames, err := client.Database("homepedia").ListCollectionNames(context.TODO(), bson.D{})
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la liste des collections: %v", err)
		}

		collectionExists := false
		for _, name := range collectionNames {
			if name == "figaroArticles" {
				collectionExists = true
				break
			}
		}

		if !collectionExists {
			err = client.Database("homepedia").CreateCollection(context.TODO(), "figaroArticles")
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la création de la collection 'figaroArticles': %v", err)
			}
			fmt.Println("Collection 'figaroArticles' créée avec succès dans la base de données existante 'homepedia'.")
		}
	}

	return client, nil
}
