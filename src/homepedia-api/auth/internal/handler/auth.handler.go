package handler

import (
	helloUsecase "homepedia-api/auth/internal/usecase"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	EchoInstance *echo.Echo
}

type AuthHandlerInterface interface {
	Register()
}

func NewAuthHandler(e *echo.Echo) AuthHandlerInterface {
	return &AuthHandler{
		EchoInstance: e,
	}
}

func (ah *AuthHandler) Register() {
	ah.EchoInstance.GET("/", helloUsecase.Execute)
}
