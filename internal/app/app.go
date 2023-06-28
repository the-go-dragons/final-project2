package app

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	customeMiddleware "github.com/the-go-dragons/final-project2/internal/app/middleware"
	handlers "github.com/the-go-dragons/final-project2/internal/interfaces/http"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"github.com/the-go-dragons/final-project2/pkg/config"
)

var (
	store = sessions.NewCookieStore()
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()
	routing(e)

	return &App{
		E: e,
	}
}

func (application *App) Start(portAddress string) error {
	err := application.E.Start(fmt.Sprintf(":%s", portAddress))
	application.E.Logger.Fatal(err)
	return err
}

func routing(e *echo.Echo) {

	initializeSessionStore()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(SessionMiddleware())

	userRepo := persistence.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.GET("/logout", userHandler.Logout, customeMiddleware.RequireAuth)
}

func initializeSessionStore() {
	secret := config.GetEnv("JWT_TOKEN_EXPIRE_HOURS")
	store = sessions.NewCookieStore([]byte(secret))

	// Set session options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // Session expiration time (in seconds)
		HttpOnly: true,
	}
}

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			session, _ := store.Get(c.Request(), "go-dragon-session")
			c.Set("session", session)

			return next(c)
		}
	}
}
