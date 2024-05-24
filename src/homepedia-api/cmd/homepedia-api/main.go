package main

import (
	"fmt"
	"homepedia-api/internal/router"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Print("Hello, World!")
	echoInstance := echo.New()

	// Router
	authRouter := router.NewAuthRouter(echoInstance)
	authRouter.Register()
	// Start server
	echoInstance.Logger.Fatal(echoInstance.Start(":1323"))
}
