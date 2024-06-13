package router

import (
	"homepedia-api/auth/internal/http/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthRouter struct {
	EchoInstance *echo.Echo
	DbInstance   *gorm.DB
}

type AuthRouterInterface interface {
	Register()
}

func NewAuthRouter(e *echo.Echo, g *gorm.DB) AuthRouterInterface {
	return &AuthRouter{
		EchoInstance: e,
		DbInstance:   g,
	}
}

func (ar *AuthRouter) Register() {
	handler := handler.NewAuthHandler(ar.EchoInstance)
	handler.Register()
}
