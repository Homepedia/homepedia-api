package auth

import (
	"homepedia-api/auth/internal/router"

	"github.com/labstack/echo/v4"
)

func InitService(echoInstance *echo.Echo) {
	authRouter := router.NewAuthRouter(echoInstance)
	authRouter.Register()
}
