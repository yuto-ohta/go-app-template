package controller

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/apputil"
	"go-app-template/src/domain"
	"go-app-template/src/domain/valueobject"
	"go-app-template/src/usecase"
	"net/http"
	"strconv"
	"strings"
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
	userId := valueobject.NewUserIdWithId(id)

	// get user
	var user domain.User
	user, err = u.userUseCase.FindById(*userId)
	var appErr *apperror.AppError
	if err != nil {
		// 該当のユーザーが存在しない場合
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, appErr, message.UserNotFound)
		}
		// 予期せぬエラー
		return apperror.ResponseErrorJSON(c, err, message.SystemError)
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
	userName, err = getUserNameParam(c.QueryParam("name"))
	if err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserName)
	}

	// new domain user
	userDomain := *domain.NewUser(userName)

	// register user domain
	var user domain.User
	user, err = u.userUseCase.CreateUser(userDomain)

	var appErr *apperror.AppError
	if err != nil {
		// ユーザー登録失敗
		if errors.As(err, &appErr) {
			return apperror.ResponseErrorJSON(c, appErr, message.CreateUserFailed)
		}
		// 予期せぬエラー
		return apperror.ResponseErrorJSON(c, err, message.SystemError)
	}

	return c.JSON(http.StatusOK, user)
}

func getUserIdParam(param string) (int, error) {
	id, err := strconv.Atoi(param)

	// 数字以外はNG
	if err != nil {
		appErr := apperror.NewAppError(err, http.StatusBadRequest)
		return 0, appErr
	}

	return id, nil
}

func getUserNameParam(param string) (string, error) {
	// 両端の半角・全角スペースを取り除く
	param = strings.TrimSpace(param)

	// 空文字はNG
	if param == "" {
		appErr := apperror.NewAppError(fmt.Errorf("\"userName\"が空文字になっています"), http.StatusBadRequest)
		return "", appErr
	}

	// 半角・全角スペース, 改行を含む場合はNG
	if apputil.ContainsSpaces(param) {
		appErr := apperror.NewAppError(fmt.Errorf("\"userName\"に半角・全角スペース, 改行コードが含まれています"), http.StatusBadRequest)
		return "", appErr
	}

	return param, nil
}
