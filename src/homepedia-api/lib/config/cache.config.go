package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitCache() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DRAGON_FLY_HOST"),   // ou l'adresse de votre serveur Redis
		Password: os.Getenv("DRAGON_FLY_SECRET"), // Mot de passe, s'il est configuré
		DB:       0,                              // Base de données à sélectionner
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
}

func GetCache() *redis.Client {
	return RedisClient
}

func CloseCache() {
	_ = RedisClient.Close()
}


