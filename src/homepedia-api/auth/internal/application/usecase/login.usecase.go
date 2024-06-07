package usecase

import (
	"fmt"
	"homepedia-api/auth/internal/application/dto"
	"homepedia-api/lib/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginExecute(c echo.Context) error {
	var req dto.UserLoginDTO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}
	fmt.Println(req)
	return nil
}
