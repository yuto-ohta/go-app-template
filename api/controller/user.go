package controller

import (
	"github.com/labstack/echo/v4"
	"go-app-template/domain"
	"go-app-template/usecase"
	"gorm.io/gorm"
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userId := domain.NewUserId(id)

	// get user
	var user domain.User
	user, err = u.userUseCase.FindById(*userId)
	
	if err == gorm.ErrRecordNotFound {
		// RecordNotFoundのときは404を返す
		return c.JSON(http.StatusNotFound, err.Error())
	} else if err != nil {
		// それ以外の意図しないエラーが返ったときは500を返す
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func getUserIdParam(c echo.Context) (int, error) {
	var id int
	var err error

	// check is "id" not empty
	if c.Param("id") == "" {
		err = fmt.Errorf("missing required argument: \"id\"")
		return id, err
	}

	// check is "id" string
	if reflect.TypeOf(c.Param("id")).Kind() != reflect.String {
		err = fmt.Errorf("\"id\" must be string: %v", c.Param("id"))
		return id, err
	}

	// parse string to int
	id, err = strconv.Atoi(c.Param("id"))

	// check parsing error has occurred
	if err != nil {
		return id, err
	}

	return id, nil
}
