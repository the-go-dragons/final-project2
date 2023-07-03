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

	walletRepo := persistence.NewWalletRepository()
	subscrptionRepo := persistence.NewSubscriptionRepository()

	userRepo := persistence.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo, walletRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	paymentRepo := persistence.NewPaymentRepository()
	paymentService := usecase.NewPayment(paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	trxRepo := persistence.NewTransactionRepository()
	walletService := usecase.NewWallet(walletRepo, paymentRepo, trxRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	numberRepo := persistence.NewNumberRepository()
	numberService := usecase.NewNumber(numberRepo, walletRepo, subscrptionRepo)
	numberHandler := handlers.NewNumberHandler(numberService, walletService)

	phonebookRepo := persistence.NewPhoneBookRepository()
	phoneBookService := usecase.NewPhoneBook(phonebookRepo, userRepo)
	phoneBookHandler := handlers.NewPhoneBookHandler(phoneBookService)

	contactRepo := persistence.NewContactRepository()
	contactService := usecase.NewContact(phonebookRepo, contactRepo, numberRepo, subscrptionRepo)
	contactHandler := handlers.NewContactHandler(contactService)

	smsRepository := persistence.NewSmsHistoryRepository()
	smsService := usecase.NewSmsService(smsRepository, *userRepo, phonebookRepo, numberRepo, subscrptionRepo, contactRepo)
	smsHandler := handlers.NewSmsHandler(smsService, contactService)

	smsTemplateRepo := persistence.NewSmsTemplateRepository()
	smsTemplateUsecase := usecase.NewSmsTemplateUsecase(smsTemplateRepo)
	smsTemplateHandler := handlers.NewSmsTemplateHandler(smsTemplateUsecase)

	adminHandler := handlers.NewAdminHandler(userUsecase)

	// TODO: add /users route prefix
	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.GET("/logout", userHandler.Logout, customeMiddleware.RequireAuth)

	e.GET("/payments/pay/:paymentId", paymentHandler.Pay)
	e.POST("/payments/callback", paymentHandler.Callback)

	e.POST("/wallets/charge-request", walletHandler.CharageRequest)
	e.POST("/wallets/finalize-charge", walletHandler.FinalizeCharge)

	e.PUT("/numbers", numberHandler.Create)
	e.POST("/numbers/buy-rent", numberHandler.BuyOrRent, customeMiddleware.RequireAuth)

	e.GET("/phonebook", phoneBookHandler.GetAll)
	e.GET("/phonebook/username", phoneBookHandler.GetByUserName)
	e.DELETE("/phonebook", phoneBookHandler.Delete)
	e.POST("/phonebook", phoneBookHandler.Create, customeMiddleware.RequireAuth)
	e.PUT("/phonebook", phoneBookHandler.Edit, customeMiddleware.RequireAuth)

	e.POST("/contact", contactHandler.Create)
	e.PUT("/contact", contactHandler.Edit)
	e.GET("/contact", contactHandler.GetAll)
	e.GET("/contact/phonebook", contactHandler.GetByPhoneBook)
	e.DELETE("/contact", contactHandler.Delete)

	e.POST("/sms", smsHandler.SendSMS, customeMiddleware.RequireAuth)
	e.POST("/sms/username", smsHandler.SendSMSByUsername, customeMiddleware.RequireAuth)

	e.POST("/templates/new", smsTemplateHandler.NewSmsTemplate, customeMiddleware.RequireAuth)

	e.GET("/admin/disable-user/:userId", adminHandler.DisableUser, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
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
