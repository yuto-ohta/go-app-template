package route

import (
	"go-app-template/src/api/controller"
	"go-app-template/src/config"
	"go-app-template/src/infrastructure"
	"go-app-template/src/usecase/impl"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _cookieEncryptKey = config.GetConfig()["cookie_encrypt_key"].(string)

func NewRouter() *echo.Echo {
	// setup
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(_cookieEncryptKey))))

	// make instance
	// -------- user
	userRepository := infrastructure.NewUserRepositoryImpl()
	userUseCase := impl.NewUserUseCaseImpl(userRepository)
	userController := controller.NewUserController(userUseCase)
	// -------- app
	appUseCase := impl.NewApplicationUseCaseImpl(userRepository)
	appController := controller.NewApplicationController(appUseCase)

	// routing
	// -------- accessible without Login
	e.GET("/users/:id", userController.GetUser)
	e.GET("/users", userController.GetAllUser)
	e.POST("/users/new", userController.CreateUser)
	e.POST("/login", appController.Login)
	e.GET("/logout", appController.Logout)
	// -------- NOT accessible without Login
	e.PUT("/users/:id/update", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	return e
}
