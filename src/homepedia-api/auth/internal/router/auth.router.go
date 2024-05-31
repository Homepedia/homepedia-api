package router

import (
	"homepedia-api/auth/internal/handler"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	EchoInstance *echo.Echo
}

type AuthRouterInterface interface {
	Register()
}

func NewAuthRouter(e *echo.Echo) AuthRouterInterface {
	return &AuthRouter{
		EchoInstance: e,
	}
}

func (ar *AuthRouter) Register() {
	handler := handler.NewAuthHandler(ar.EchoInstance)
	handler.Register()
}
