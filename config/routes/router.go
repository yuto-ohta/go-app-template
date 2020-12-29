package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-app-template/api/controller"
	"go-app-template/infrastructure"
	"go-app-template/usecase/impl"
)

func NewRouter() *echo.Echo {
	// setup
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// make instance
	userRepository := infrastructure.NewUserRepositoryImpl()
	userUseCase := impl.NewUserUseCaseImpl(userRepository)
	userController := controller.NewUserController(userUseCase)

	// routing
	e.GET("/user/:id", userController.GetUser)

	return e
}
