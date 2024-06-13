package handler

import (
	"homepedia-api/auth/internal/application/usecase"

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
	ah.EchoInstance.POST("/auth/register", usecase.RegisterExecute)
	ah.EchoInstance.POST("/auth/login", usecase.LoginExecute)
}
