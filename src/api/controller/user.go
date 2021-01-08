package controller

import (
	"errors"
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: useCase}
}

/*
	ユーザーを取得する
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
	var user dto.UserDto
	if user, err = u.userUseCase.FindById(id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, appErr, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.GetUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/*
	ユーザーを全件取得する
	@return users
*/
func (u UserController) GetAll(c echo.Context) error {
	var err error

	// get all user
	var users []dto.UserDto
	if users, err = u.userUseCase.FindAll(); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.GetUserFailed)
	}

	return c.JSON(http.StatusOK, users)
}

/*
	ユーザーを新規登録する
	@body_param name: userName
	@return user
*/
func (u UserController) CreateUser(c echo.Context) error {
	var err error

	// get userDto
	var userDto dto.UserDto
	if err = c.Bind(&userDto); err != nil {
		appErr := apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		return apperror.ResponseErrorJSON(c, appErr, message.StatusBadRequest)
	}
	if err = userDto.Validate(); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// register user
	var user dto.UserDto
	if user, err = u.userUseCase.CreateUser(userDto.Name); err != nil {
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
	var user dto.UserDto
	if user, err = u.userUseCase.DeleteUser(id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.DeleteUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/*
	ユーザーを更新する
	@path_param id: userId
	@body_param name: userName
	@return user
*/
func (u UserController) UpdateUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// get userDto
	var newUser dto.UserDto
	if err = c.Bind(&newUser); err != nil {
		appErr := apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		return apperror.ResponseErrorJSON(c, appErr, message.StatusBadRequest)
	}
	if err = newUser.Validate(); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// update user
	var updated dto.UserDto
	if updated, err = u.userUseCase.UpdateUser(id, newUser); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.UpdateUserFailed)
	}

	return c.JSON(http.StatusOK, updated)
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
