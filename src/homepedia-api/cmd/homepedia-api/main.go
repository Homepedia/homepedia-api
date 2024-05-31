package main

import (
	auth "homepedia-api/auth"
	"homepedia-api/lib/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()
	echoInstance := echo.New()
	// Init connections
	config.Init()

	// Init auth service
	auth.InitService(echoInstance)

	// Start server
	echoInstance.Logger.Fatal(echoInstance.Start(":1323"))
}
