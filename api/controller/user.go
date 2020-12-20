package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-app-template/domain"
	"go-app-template/usecase"
	"net/http"
	"reflect"
	"strconv"
)

type UserController struct {
}

/*
	@path_param userId
	@return user
*/
func GetUser(c echo.Context) error {
	var err error

	// get userId
	var id int
	if id, err = getUserIdParam(c); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// get user
	var user domain.User
	userUseCase := usecase.UserUseCase{}
	if user, err = userUseCase.FindById(id); err != nil {
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
