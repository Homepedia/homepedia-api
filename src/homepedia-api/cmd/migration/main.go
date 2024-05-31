package main

import (
	"homepedia-api/lib/config"
	auth_domain "homepedia-api/lib/domain/auth"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Init connection to databases
	config.Init()

	// Get auth database connection
	authDB := config.Connections.Auth

	// Downgrade auth database
	authDB.Migrator().DropTable(&auth_domain.Credentials{})
	authDB.Migrator().DropTable(&auth_domain.Role{})

	// Migrate auth database
	authDB.AutoMigrate(&auth_domain.Credentials{})
	authDB.AutoMigrate(&auth_domain.Role{})

	// Generate seeds for roles
	roles := []auth_domain.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	authDB.CreateInBatches(roles, len(roles))

}
