package main

import (
	"homepedia-api/lib/config"
	"homepedia-api/lib/domain"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Init connection to databases
	config.Init()

	// Get auth database connection
	authDB := config.Connections.Auth

	// Downgrade auth database
	authDB.Migrator().DropTable(&domain.Credentials{})
	authDB.Migrator().DropTable(&domain.Role{})

	// Migrate auth database
	authDB.AutoMigrate(&domain.Credentials{})
	authDB.AutoMigrate(&domain.Role{})

	// Generate seeds for roles
	roles := []domain.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	authDB.CreateInBatches(roles, len(roles))

}
