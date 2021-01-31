package usecase

import (
	"go-app-template/src/api/controller/dto"

	"github.com/labstack/echo/v4"
)

type AuthenticationUseCase interface {
	Authenticate(c echo.Context, targetUserIdInt int) (dto.UserResDto, error)
}
