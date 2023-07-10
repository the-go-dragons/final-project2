package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	customeMiddleware "github.com/the-go-dragons/final-project2/internal/app/middleware"
	handlers "github.com/the-go-dragons/final-project2/internal/interfaces/http"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"github.com/the-go-dragons/final-project2/internal/usecase"
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

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
	contactHandler := handlers.NewContactHandler(contactService, phoneBookService)

	smsRepository := persistence.NewSmsHistoryRepository()
	smsService := usecase.NewSmsService(smsRepository, *userRepo, phonebookRepo, numberRepo, subscrptionRepo, contactRepo)
	smsHandler := handlers.NewSmsHandler(smsService, contactService, phoneBookService)

	smsTemplateRepo := persistence.NewSmsTemplateRepository()
	smsTemplateUsecase := usecase.NewSmsTemplateUsecase(smsTemplateRepo)
	smsTemplateHandler := handlers.NewSmsTemplateHandler(smsTemplateUsecase, smsService, contactService, phoneBookService)

	priceRepo := persistence.NewPriceRepository()
	priceUsecase := usecase.NewPriceService(priceRepo)

	adminHandler := handlers.NewAdminHandler(*userUsecase, priceUsecase, smsService)

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.GET("/logout", userHandler.Logout, customeMiddleware.RequireAuth)

	e.GET("/payments/pay/:paymentId", paymentHandler.Pay, customeMiddleware.RequireAuth)
	e.POST("/payments/callback", paymentHandler.Callback, customeMiddleware.RequireAuth)

	e.POST("/wallets/charge-request", walletHandler.CharageRequest, customeMiddleware.RequireAuth)
	e.POST("/wallets/finalize-charge", walletHandler.FinalizeCharge, customeMiddleware.RequireAuth)

	e.PUT("/numbers", numberHandler.Create, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.POST("/numbers/buy-rent", numberHandler.BuyOrRent, customeMiddleware.RequireAuth)

	e.GET("/phonebook", phoneBookHandler.GetAll, customeMiddleware.RequireAuth)
	e.GET("/phonebook/username", phoneBookHandler.GetByUserName, customeMiddleware.RequireAuth)
	e.DELETE("/phonebook", phoneBookHandler.Delete, customeMiddleware.RequireAuth)
	e.POST("/phonebook", phoneBookHandler.Create, customeMiddleware.RequireAuth)
	e.PUT("/phonebook", phoneBookHandler.Edit, customeMiddleware.RequireAuth)

	e.POST("/contact", contactHandler.Create, customeMiddleware.RequireAuth)
	e.PUT("/contact", contactHandler.Edit, customeMiddleware.RequireAuth)
	e.GET("/contact", contactHandler.GetAll, customeMiddleware.RequireAuth)
	e.GET("/contact/phonebook", contactHandler.GetByPhoneBook, customeMiddleware.RequireAuth)
	e.DELETE("/contact", contactHandler.Delete, customeMiddleware.RequireAuth)

	e.POST("/sms", smsHandler.SendSingleSMS, customeMiddleware.RequireAuth)
	e.POST("/sms/periodic", smsHandler.SendSinglePeriodSMS, customeMiddleware.RequireAuth)
	e.POST("/sms/username", smsHandler.SendSingleSMSByUsername, customeMiddleware.RequireAuth)
	e.POST("/sms/username/periodic", smsHandler.SendSinglePeriodSMSByUsername, customeMiddleware.RequireAuth)

	e.POST("/templates/new", smsTemplateHandler.NewSmsTemplate, customeMiddleware.RequireAuth)
	e.GET("/templates", smsTemplateHandler.SmsTemplateList, customeMiddleware.RequireAuth)
	e.POST("/templates/sms", smsTemplateHandler.NewSingleSmsWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/periodic", smsTemplateHandler.NewSinglePeriodSmsWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/username", smsTemplateHandler.NewSingleSmsWithUsernameWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/username/periodic", smsTemplateHandler.NewSinglePeriodSmsWithUsernameWithTemplate, customeMiddleware.RequireAuth)

	e.GET("/admin/disable-user/:userId", adminHandler.DisableUser, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/admin/change-priceing", adminHandler.ChangePricing, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/admin/sms-report/:userId", adminHandler.GetSMSHistoryByUserId, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
}
