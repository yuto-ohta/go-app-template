package route

import (
	"go-app-template/src/api/controller"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase/impl"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.GET("/users/:id", userController.GetUser)
	e.GET("/users", userController.GetAll)
	e.POST("/users/new", userController.CreateUser)
	e.DELETE("/users/:id", userController.DeleteUser)
	e.PUT("/users/:id/update", userController.UpdateUser)

	return e
}
