package usecase

import (
	"fmt"
	"homepedia-api/auth/internal/application/dto"
	"homepedia-api/lib/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Execute(c echo.Context) error {

	var req dto.UserRegisterDTO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	credentials := domain.NewCredentials(req.Username, req.Password, req.Email, 2)

	fmt.Println(credentials)

	return c.String(http.StatusOK, "Hello, World! Please stand up")
}
