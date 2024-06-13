package middleware

import (
	"encoding/json"
	"homepedia-api/lib/config/cache"
	"homepedia-api/lib/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthGuard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token, err := ctx.Cookie("session_id")
			if err != nil {
				println(err.Error())
				return ctx.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "Unauthorized"})
			}
			userId, err := cache.GetSessionId(token.Value, ctx)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "Unauthorized"})
			}

			var userTokenValue utils.UserTokenValue
			err = json.Unmarshal([]byte(userId), &userTokenValue)

			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "Unauthorized"})
			}

			ctx.Set("user", userTokenValue)
			return next(ctx)
		}
	}
}
