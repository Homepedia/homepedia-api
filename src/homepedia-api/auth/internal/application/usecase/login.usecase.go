package usecase

import (
	"errors"
	"homepedia-api/auth/internal/application/dto"
	"homepedia-api/auth/internal/http/repository"
	"homepedia-api/lib/service"
	"homepedia-api/lib/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	USER_LOGGED      = "user logged successfully"
	INVALID_PASSWORD = "invalid password"
	USER_NOT_FOUND   = "user not found"
)

func LoginExecute(c echo.Context) error {
	var req dto.UserLoginDTO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
	}

	authRepo := repository.NewAuthRepository()

	credentials, err := authRepo.FindUserByEmail(req.Email)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, utils.HttpResponse{Message: USER_NOT_FOUND})
		}
		return c.JSON(http.StatusInternalServerError, utils.HttpResponse{Message: err.Error()})
	}

	if !service.VerifyPassword(req.Password, credentials.Password) {
		return c.JSON(http.StatusForbidden, utils.HttpResponse{Message: INVALID_PASSWORD})
	}

	return c.JSON(http.StatusOK, utils.HttpResponse{Message: USER_LOGGED})
}
