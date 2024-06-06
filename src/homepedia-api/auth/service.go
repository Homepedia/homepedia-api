package auth

import (
	"homepedia-api/auth/internal/http/router"
	"homepedia-api/lib/config"

	"github.com/labstack/echo/v4"
)

func InitService(echoInstance *echo.Echo) {
	dbInstance := config.Connections.Auth
	authRouter := router.NewAuthRouter(echoInstance, dbInstance)
	authRouter.Register()
}
