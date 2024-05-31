package auth

import (
	"fmt"
	"homepedia-api/auth/internal/router"
	"homepedia-api/lib/config"

	"github.com/labstack/echo/v4"
)

func InitService(echoInstance *echo.Echo) {
	authRouter := router.NewAuthRouter(echoInstance)
	authRouter.Register()
	db := config.Connections.Auth
	fmt.Println(db)
}
