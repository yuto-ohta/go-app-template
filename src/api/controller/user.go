package controller

import (
	"errors"
	"fmt"
	"go-app-template/src/api/controller/dto"
	"go-app-template/src/apperror"
	"go-app-template/src/apperror/message"
	"go-app-template/src/usecase"
	"go-app-template/src/usecase/appmodel"
	"go-app-template/src/usecase/appmodel/session"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUseCase usecase.UserUseCase
	authUseCase usecase.AuthenticationUseCase
}

func NewUserController(useCase usecase.UserUseCase, authUseCase usecase.AuthenticationUseCase) *UserController {
	return &UserController{
		userUseCase: useCase,
		authUseCase: authUseCase,
	}
}

/**************************************
	ユーザーをuserIdで取得する
	@path_param id: userId
	@return user
**************************************/

func (u UserController) GetUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// get user
	var user dto.UserResDto
	if user, err = u.userUseCase.GetUser(id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, appErr, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.GetUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/**************************************
	ユーザーを全件取得する(オプション: orderBy, order, page, limit)
	@query_param orderBy, order, page, limit
	@return users
**************************************/

func (u UserController) GetAllUser(c echo.Context) error {
	var err error

	// get orderBy
	var orderBy dto.UserSortColumn
	if orderBy, err = getOptionalOrderByParam(c.QueryParam("orderBy")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// get order
	var order appmodel.Order
	if order, err = getOptionalOrderParam(c.QueryParam("order")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// get page
	var page int
	if page, err = getOptionalQueryParamInt(c.QueryParam("page")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// get limit
	var limit int
	if limit, err = getOptionalQueryParamInt(c.QueryParam("limit")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// make searchCondition
	var condition *appmodel.SearchCondition
	if condition, err = appmodel.NewSearchCondition(orderBy.String(), order, page, limit); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// get users
	var userPage dto.UserPage
	if userPage, err = u.userUseCase.GetAllUser(*condition); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.GetUserFailed)
	}

	return c.JSON(http.StatusOK, userPage)
}

/**************************************
	ユーザーを新規登録する
	@body_param name: userName, password: password
	@return user
**************************************/

func (u UserController) CreateUser(c echo.Context) error {
	var err error

	// get userDto
	var userDto dto.UserReceiveDto
	if err = c.Bind(&userDto); err != nil {
		appErr := apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
		return apperror.ResponseErrorJSON(c, appErr, message.StatusBadRequest)
	}
	if err = userDto.Validate(); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// register user
	var user dto.UserResDto
	if user, err = u.userUseCase.CreateUser(userDto); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.CreateUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/**************************************
	ユーザーを削除する
	@path_param id: userId
	@return user
**************************************/

func (u UserController) DeleteUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// authenticate
	if _, err = u.authUseCase.Authenticate(c, id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusUnauthorized {
			return apperror.ResponseErrorJSON(c, err, message.UnAuthorized)
		}
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusForbidden {
			return apperror.ResponseErrorJSON(c, err, message.Forbidden)
		}
		return apperror.ResponseErrorJSON(c, err, message.DeleteUserFailed)
	}

	// delete user
	var user dto.UserResDto
	if user, err = u.userUseCase.DeleteUser(id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.DeleteUserFailed)
	}

	return c.JSON(http.StatusOK, user)
}

/**************************************
	ユーザーを更新する
	@path_param id: userId
	@body_param name: userName, password: password
	@return user
**************************************/

func (u UserController) UpdateUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c.Param("id")); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.InvalidUserId)
	}

	// get userDto
	var newUser dto.UserReceiveDto
	if newUser, err = getUpdateUserBodyParam(c); err != nil {
		return apperror.ResponseErrorJSON(c, err, message.StatusBadRequest)
	}

	// authenticate
	if _, err = u.authUseCase.Authenticate(c, id); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusUnauthorized {
			return apperror.ResponseErrorJSON(c, err, message.UnAuthorized)
		}
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusForbidden {
			return apperror.ResponseErrorJSON(c, err, message.Forbidden)
		}
		return apperror.ResponseErrorJSON(c, err, message.UpdateUserFailed)
	}

	// update user
	var updated dto.UserResDto
	if updated, err = u.userUseCase.UpdateUser(id, newUser); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.GetHttpStatus() == http.StatusNotFound {
			return apperror.ResponseErrorJSON(c, err, message.UserNotFound)
		}
		return apperror.ResponseErrorJSON(c, err, message.UpdateUserFailed)
	}

	// make user to login after change password
	if len(newUser.Password) != 0 {
		// invalidate session
		manager := session.NewSessionManager()
		if err = manager.InvalidateSession(c); err != nil {
			return apperror.ResponseErrorJSON(c, err, message.LogoutFailed)
		}
	}

	return c.JSON(http.StatusOK, updated)
}

/**************************************
	private
**************************************/

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

func getOptionalQueryParamInt(param string) (int, error) {
	var err error

	// 未指定はOK
	if len(param) == 0 {
		return -1, nil
	}

	// 数字以外はNG
	var paramInt int
	if paramInt, err = strconv.Atoi(param); err != nil {
		return -1, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	// 0以下はNG
	if paramInt <= 0 {
		return -1, apperror.NewAppErrorWithStatus(fmt.Errorf("対象のクエリパラメータは1以上を指定してください"), http.StatusBadRequest)
	}

	return paramInt, nil
}

func getOptionalOrderByParam(param string) (dto.UserSortColumn, error) {
	// 未指定の場合は、ソートColumnはIDとする
	if len(param) == 0 {
		return dto.ID, nil
	}

	switch param {
	case dto.ID.String():
		return dto.ID, nil
	case dto.USERNAME.String():
		return dto.USERNAME, nil
	default:
		return -1, apperror.NewAppErrorWithStatus(fmt.Errorf("指定のColumnはソートに使用できるものではありません, param: %v", param), http.StatusBadRequest)
	}
}

func getOptionalOrderParam(param string) (appmodel.Order, error) {
	// 未指定の場合は、orderはASCとする
	if len(param) == 0 {
		return appmodel.ASC, nil
	}

	order := appmodel.GetOrder(param)
	if order < 0 {
		return -1, apperror.NewAppErrorWithStatus(fmt.Errorf("orderには\"ASC\",\"DESC\"のどちらかを使用してください, param: %v", param), http.StatusBadRequest)
	}

	return order, nil
}

func getUpdateUserBodyParam(c echo.Context) (dto.UserReceiveDto, error) {
	var err error
	var receiveDto dto.UserReceiveDto
	if err = c.Bind(&receiveDto); err != nil {
		return dto.UserReceiveDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	// ユーザー名 && パスワード が両方取得できない場合はNG
	if len(receiveDto.Name) == 0 && len(receiveDto.Password) == 0 {
		appErr := apperror.NewAppErrorWithStatus(fmt.Errorf("ユーザー名・パスワードが取得できません"), http.StatusBadRequest)
		return dto.UserReceiveDto{}, appErr
	}
	// ユーザー名がある場合のValidation
	if err = dto.ValidateName(receiveDto.Name); err != nil && len(receiveDto.Name) != 0 {
		return dto.UserReceiveDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	// パスワードがある場合のValidation
	if err = dto.ValidatePassword(receiveDto.Password); err != nil && len(receiveDto.Password) != 0 {
		return dto.UserReceiveDto{}, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}
	return receiveDto, nil
}
