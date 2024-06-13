package router

import (
	"homepedia-api/auth/internal/http/handler"
	"homepedia-api/auth/internal/http/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TestRouter struct {
    EchoInstance *echo.Echo
    DbInstance   *gorm.DB
}

type TestRouterInterface interface {
    Register()
}

func NewTestRouter(e *echo.Echo, g *gorm.DB) TestRouterInterface {
    return &TestRouter{
        EchoInstance: e,
        DbInstance:   g,
    }
}

func (ar *TestRouter) Register() {
    g := ar.EchoInstance.Group("/test")
    g.Use(middleware.AuthGuard())

    handler := handler.NewTestHandler(ar.EchoInstance)
    handler.Register(g)
}
