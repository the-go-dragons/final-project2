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
	store     = sessions.NewCookieStore()
	getSecret = func() string {
		return config.Config.Jwt.Token.Secret.Key
	}
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

func (application *App) Start(portAddress int) error {
	err := application.E.Start(fmt.Sprintf(":%d", portAddress))
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

	paymentRepo := persistence.NewPaymentRepository()
	paymentService := usecase.NewPayment(paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	walletRepo := persistence.NewWalletRepository()
	trxRepo := persistence.NewTransactionRepository()
	walletService := usecase.NewWallet(walletRepo, paymentRepo, trxRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	smsTemplateRepo := persistence.NewSmsTemplateRepository()
	smsTemplateUsecase := usecase.NewSmsTemplateUsecase(smsTemplateRepo)
	smsTemplateHandler := handlers.NewSmsTemplateHandler(smsTemplateUsecase)

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.GET("/logout", userHandler.Logout, customeMiddleware.RequireAuth)
	e.GET("/payments/pay/:paymentId", paymentHandler.Pay)
	e.POST("/payments/callback", paymentHandler.Callback)
	e.POST("/wallets/charge-request", walletHandler.CharageRequest)
	e.POST("/wallets/finalize-charge", walletHandler.FinalizeCharge)
	e.POST("/new-templates/new", smsTemplateHandler.NewSmsTemplate, customeMiddleware.RequireAuth)
}

func initializeSessionStore() {
	store = sessions.NewCookieStore([]byte(getSecret()))

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
