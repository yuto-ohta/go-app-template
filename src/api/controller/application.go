package controller

import (
	"errors"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/usecase"
	"go-app-template/src/usecase/appmodel"
	"go-app-template/src/usecase/appmodel/session"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApplicationController struct {
	applicationUseCase usecase.ApplicationUseCase
}

func NewApplicationController(useCase usecase.ApplicationUseCase) *ApplicationController {
	return &ApplicationController{applicationUseCase: useCase}
}

/**************************************
	ログイン
	@body_param id: userId, password: password
	@return NoContent
**************************************/

func (a ApplicationController) Login(c echo.Context) error {
	var err error

	// get loginDto
	var loginDto dto.LoginReceiveDto
	if err = c.Bind(&loginDto); err != nil {
		appErr := apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		return apperror.ResponseErrorJSON(c, appErr, message.StatusBadRequest)
	}
	if err = loginDto.Validate(); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// try login
	var token appmodel.SignedToken
	if token, err = a.applicationUseCase.Login(loginDto); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusUnauthorized {
			return apperror.ResponseErrorJSON(c, err, message.WrongPassword)
		}
		return apperror.ResponseErrorJSON(c, err, message.LoginFailed)
	}

	// set token into session
	manager := session.NewSessionManager()
	if err = manager.SetSession(c, "token", string(token)); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.LoginFailed)
	}

	return c.NoContent(http.StatusOK)
}

/**************************************
	ログアウト
	@return NoContent
**************************************/

func (a ApplicationController) Logout(c echo.Context) error {
	var err error

	// invalidate session
	manager := session.NewSessionManager()
	if err = manager.InvalidateSession(c); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.LogoutFailed)
	}

	return c.NoContent(http.StatusOK)
}
