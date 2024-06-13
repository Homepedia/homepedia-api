package cache

import (
	"encoding/json"
	"homepedia-api/lib/config"
	"homepedia-api/lib/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func GetSessionId(sessionId string, ctx echo.Context) (string, error) {
	return config.GetCache().Get(ctx.Request().Context(), sessionId).Result()
}

func SetSessionId(sessionId string, value interface{}, ctx echo.Context) error {
	return config.GetCache().Set(ctx.Request().Context(), sessionId, value, 6*time.Hour).Err()
}

func CreateSessionId(payload interface{}, ctx echo.Context) (string, error) {
	sessionId := utils.GenerateSessionID()
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	err = SetSessionId(sessionId, data, ctx)
	return sessionId, err
}

func DeleteSessionId(sessionId string, ctx echo.Context) error {
	return config.GetCache().Del(ctx.Request().Context(), sessionId).Err()
}
