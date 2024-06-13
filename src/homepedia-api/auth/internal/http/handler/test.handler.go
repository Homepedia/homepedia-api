package handler

import (
	"homepedia-api/lib/utils"

	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	EchoInstance *echo.Echo
}

type TestHandlerInterface interface {
	Register(*echo.Group)
}

func NewTestHandler(e *echo.Echo) TestHandlerInterface {
	return &TestHandler{
		EchoInstance: e,
	}
}

func (th *TestHandler) Register(g *echo.Group) {
	g.GET("", test)
}

func test(ctx echo.Context) error {
	userTokenValue := ctx.Get("user").(utils.UserTokenValue)
	return ctx.JSON(200, userTokenValue.UserId)
}
