package app

import (
	"fmt"
	"net/http"

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

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	walletRepo := persistence.NewWalletRepository()
	subscrptionRepo := persistence.NewSubscriptionRepository()

	paymentRepo := persistence.NewPaymentRepository()
	paymentService := usecase.NewPayment(paymentRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	trxRepo := persistence.NewTransactionRepository()
	walletService := usecase.NewWallet(walletRepo, paymentRepo, trxRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	numberRepo := persistence.NewNumberRepository()
	numberService := usecase.NewNumber(numberRepo, walletRepo, subscrptionRepo)
	numberHandler := handlers.NewNumberHandler(numberService, walletService)

	userRepo := persistence.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo, walletRepo, numberRepo, subscrptionRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	phonebookRepo := persistence.NewPhoneBookRepository()
	phoneBookService := usecase.NewPhoneBook(phonebookRepo, userRepo)
	phoneBookHandler := handlers.NewPhoneBookHandler(phoneBookService)

	contactRepo := persistence.NewContactRepository()
	contactService := usecase.NewContact(phonebookRepo, contactRepo, numberRepo, subscrptionRepo)
	contactHandler := handlers.NewContactHandler(contactService, phoneBookService)

	wordRepo := persistence.NewInappropriateWordRepository()
	wordService := usecase.NewInappropriateWord(wordRepo)
	wordHandler := handlers.NewInappropriateWordHandler(wordService)

	smsRepository := persistence.NewSmsHistoryRepository()
	smsService := usecase.NewSmsService(smsRepository, userRepo, phonebookRepo, numberRepo, subscrptionRepo, contactRepo)
	smsHandler := handlers.NewSmsHandler(smsService, contactService, phoneBookService, wordService)

	smsTemplateRepo := persistence.NewSmsTemplateRepository()
	smsTemplateUsecase := usecase.NewSmsTemplateService(smsTemplateRepo)
	smsTemplateHandler := handlers.NewSmsTemplateHandler(smsTemplateUsecase, smsService, contactService, phoneBookService)

	priceRepo := persistence.NewPriceRepository()
	priceUsecase := usecase.NewPriceService(priceRepo)

	adminHandler := handlers.NewAdminHandler(userUsecase, priceUsecase, smsService)

	smsHistoryrepo := persistence.NewSmsHistoryRepository()
	smsHistoryUsecase := usecase.NewSmsHistoryUsecase(smsHistoryrepo)
	smsHistoryHandler := handlers.NewSmsHistoryHandler(smsHistoryUsecase)

	e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, "welcome to sms panel Q")
    })

	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)
	e.GET("/logout", userHandler.Logout, customeMiddleware.RequireAuth)
	e.POST("/users/:userId/update-default-number", userHandler.UpdateDefaultNumber)

	e.GET("/payments/pay/:paymentId", paymentHandler.Pay, customeMiddleware.RequireAuth)
	e.POST("/payments/callback", paymentHandler.Callback)

	e.POST("/wallets/charge-request", walletHandler.CharageRequest, customeMiddleware.RequireAuth)
	e.POST("/wallets/finalize-charge", walletHandler.FinalizeCharge, customeMiddleware.RequireAuth)

	e.GET("/numbers", numberHandler.GetAvailables)
	e.POST("/numbers", numberHandler.Create, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.POST("/numbers/buy-rent", numberHandler.BuyOrRent, customeMiddleware.RequireAuth)

	e.GET("/phonebook", phoneBookHandler.GetAll, customeMiddleware.RequireAuth)
	e.DELETE("/phonebook/:id", phoneBookHandler.Delete, customeMiddleware.RequireAuth)
	e.POST("/phonebook", phoneBookHandler.Create, customeMiddleware.RequireAuth)

	e.POST("/contact/:phonebookId", contactHandler.CreateContact, customeMiddleware.RequireAuth)
	e.GET("/contact/:phonebookId", contactHandler.GetByPhoneBook, customeMiddleware.RequireAuth)
	e.DELETE("/contact/:phonebookId", contactHandler.DeleteContact, customeMiddleware.RequireAuth)

	e.POST("/sms", smsHandler.SendSingleSMS, customeMiddleware.RequireAuth)
	e.POST("/sms/periodic", smsHandler.SendSinglePeriodSMS, customeMiddleware.RequireAuth)
	e.POST("/sms/username", smsHandler.SendSingleSMSByUsername, customeMiddleware.RequireAuth)
	e.POST("/sms/username/periodic", smsHandler.SendSinglePeriodSMSByUsername, customeMiddleware.RequireAuth)
	e.POST("/sms/phonebooks", smsHandler.SendSMSToPhonebooks, customeMiddleware.RequireAuth)
	e.POST("/sms/phonebooks/periodic", smsHandler.SendPeriodSMSToPhonebooks, customeMiddleware.RequireAuth)

	e.POST("/templates/new", smsTemplateHandler.NewSmsTemplate, customeMiddleware.RequireAuth)
	e.GET("/templates", smsTemplateHandler.SmsTemplateList, customeMiddleware.RequireAuth)
	e.POST("/templates/sms", smsTemplateHandler.NewSingleSmsWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/periodic", smsTemplateHandler.NewSinglePeriodSmsWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/username", smsTemplateHandler.NewSingleSmsWithUsernameWithTemplate, customeMiddleware.RequireAuth)
	e.POST("/templates/sms/username/periodic", smsTemplateHandler.NewSinglePeriodSmsWithUsernameWithTemplate, customeMiddleware.RequireAuth)

	e.GET("/admin/disable-user/:userId", adminHandler.DisableUser, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/admin/change-priceing", adminHandler.ChangePricing, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/admin/sms-report/:userId", adminHandler.GetSMSHistoryByUserId, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/admin/sms-history/search", smsHistoryHandler.Search, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)

	e.POST("/inappropriate-word", wordHandler.CreateInappropriateWord, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.GET("/inappropriate-word", wordHandler.GetAll, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
	e.DELETE("/inappropriate-word/:id", wordHandler.Delete, customeMiddleware.RequireAuth, customeMiddleware.RequireAdmin)
}
