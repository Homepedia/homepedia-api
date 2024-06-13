package usecase

import (
	"homepedia-api/auth/internal/application/dto"
	"homepedia-api/auth/internal/http/repository"
	"homepedia-api/lib/domain"
	"homepedia-api/lib/service"
	"homepedia-api/lib/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterExecute(c echo.Context) error {
	var req dto.UserRegisterDTO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}

	credentials := domain.NewCredentials(req.Username, req.Password, req.Email, 2)

	authRepo := repository.NewAuthRepository()

	hashPassword, err := service.HashPassword(credentials.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}
	credentials.Password = hashPassword

	createUser := authRepo.Register(credentials)

	if !createUser.Success {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: createUser.Message})
	}

	return c.JSON(http.StatusCreated, utils.HttpResponse{Message: createUser.Message})

}
