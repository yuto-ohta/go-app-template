package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go-app-template/domain"
	appErrors "go-app-template/errors"
	"go-app-template/errors/messages"
	"go-app-template/usecase"
	"net/http"
	"strconv"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: useCase}
}

/*
	ユーザーをuserIdで取得する
	@path_param userId
	@return user
*/
func (u UserController) GetUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c); err != nil {
		return appErrors.ResponseErrorJSON(c, err, messages.InvalidUserId.String())
	}
	userId := domain.NewUserId(id)

	// get user
	var user domain.User
	user, err = u.userUseCase.FindById(*userId)
	var appErr *appErrors.AppError
	if err != nil {
		// 該当のユーザーが存在しない場合
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return appErrors.ResponseErrorJSON(c, err, messages.UserNotFound.String())
		}
		// 予期せぬエラー
		return appErrors.ResponseErrorJSON(c, err, messages.SystemError.String())
	}

	return c.JSON(http.StatusOK, user)
}

func getUserIdParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		appErr := appErrors.NewAppError(err, http.StatusBadRequest)
		return 0, appErr
	}

	return id, nil
}
