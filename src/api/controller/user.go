package controller

import (
	"errors"
	"fmt"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/apputil"
	"go-app-template/src/domain"
	"go-app-template/src/usecase"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: useCase}
}

/*
	ユーザーをuserIdで取得する
	@path_param id: userId
	@return user
*/
func (u UserController) GetUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// get user
	var user domain.User
	if user, err = u.userUseCase.FindById(id); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

/*
	ユーザーを新規登録する
	@query_param name: userName
	@return user
*/
func (u UserController) CreateUser(c echo.Context) error {
	var err error

	// get userName
	var userName string
	if userName, err = getUserNameParam(c.QueryParam("name")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserName)
	}

	// register user
	var user domain.User
	if user, err = u.userUseCase.CreateUser(userName); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.CreateUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/*
	ユーザーを削除する
	@path_param id: userId
	@return user
*/
func (u UserController) DeleteUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// delete user
	var user domain.User
	var appErr *apperror.AppError
	if user, err = u.userUseCase.DeleteUser(id); err != nil {
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.DeleteUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

func getUserIdParam(param string) (int, error) {
	id, err := strconv.Atoi(param)

	// 数字以外はNG
	if err != nil {
		appErr := apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		return 0, appErr
	}

	// 0以下はNG
	if id <= 0 {
		appErr := apperror.NewAppErrorWithStatus(fmt.Errorf("userIdは1以上を指定してください, userId: %v", id), http.StatusBadRequest)
		return 0, appErr
	}

	return id, nil
}

func getUserNameParam(param string) (string, error) {
	// 両端の半角・全角スペースを取り除く
	param = strings.TrimSpace(param)

	// 空文字はNG
	if param == "" {
		appErr := apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"が空文字になっています, userName: %v`, param), http.StatusBadRequest)
		return "", appErr
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpace(param) {
		appErr := apperror.NewAppErrorWithStatus(fmt.Errorf(`"userName"に半角・全角スペース, 改行コードが含まれています, userName: %v`, param), http.StatusBadRequest)
		return "", appErr
	}

	// 9文字以上はNG
	if utf8.RuneCountInString(param) > 8 {
		appErr := apperror.NewAppErrorWithStatus(fmt.Errorf(`userNameは最大8文字までです, userName: %v`, param), http.StatusBadRequest)
		return "", appErr
	}

	return param, nil
}
