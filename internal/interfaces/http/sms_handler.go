package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
)

type SMSHandler interface {
	SendSingleSMS(c echo.Context) error
	SendSingleSMSByUsername(c echo.Context) error
	SendSinglePeriodSMS(c echo.Context) error
	SendSinglePeriodSMSByUsername(c echo.Context) error
	SendSMSToPhonebooks(c echo.Context) error
	SendPeriodSMSToPhonebooks(c echo.Context) error
}

type smsHandler struct {
	smsService       usecase.SMSService
	contactService   usecase.ContactService
	phoneBookService usecase.PhoneBookService
	wordService      usecase.InappropriateWordService
	priceService     usecase.PriceService
}

func NewSmsHandler(
	smsService usecase.SMSService,
	contactService usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
	wordService usecase.InappropriateWordService,
	priceService usecase.PriceService,
) SMSHandler {
	return smsHandler{
		smsService:       smsService,
		contactService:   contactService,
		phoneBookService: phoneBookService,
		wordService:      wordService,
		priceService:     priceService,
	}
}

type SingSMSRequest struct {
	SenderNumber   string `json:"senderNumber"`
	ReceiverNumber string `json:"receiverNumber"`
	Content        string `json:"content"`
}

type SingSMSWithUsernameRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverUsername string `json:"receiverUsername"`
	Content          string `json:"content"`
	PhoneBookId      uint   `json:"phoneBookId"`
}

type SingPeriodSMSRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverNumber   string `json:"receiverNumber"`
	Content          string `json:"content"`
	Period           string `json:"period"`
	RepeatationCount uint   `json:"repeatationCount"`
}

type SingPeriodSMSWithUsernameRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverUsername string `json:"receiverUsername"`
	Content          string `json:"content"`
	PhoneBookId      uint   `json:"phoneBookId"`
	Period           string `json:"period"`
	RepeatationCount uint   `json:"repeatationCount"`
}

type PhoneBookSMSRequest struct {
	SenderNumber       string `json:"senderNumber"`
	ReceiverPhoneBooks []uint `json:"receiverPhoneBooks"`
	Content            string `json:"content"`
}

type PhoneBookPeriodSMSRequest struct {
	SenderNumber       string `json:"senderNumber"`
	ReceiverPhoneBooks []uint `json:"receiverPhoneBooks"`
	Content            string `json:"content"`
	Period             string `json:"period"`
	RepeatationCount   uint   `json:"repeatationCount"`
}

func (sh smsHandler) CheckTheWalletBallence(user domain.User, receiversCount uint) (domain.Wallet, uint, error) {
	// Check the wallet and sms price
	price, err := sh.priceService.GetPrice()
	if err != nil || price.ID == 0 {
		return domain.Wallet{}, 0, errors.New("can't get price model")
	}
	wallet, err := sh.smsService.GetUserWallet(user.ID)
	if err != nil || wallet.ID == 0 {
		return domain.Wallet{}, 0, errors.New("can't get user wallet")
	}

	// Check the price type and wallet ballance
	var p uint
	if receiversCount == 1 {
		p = price.SingleSMS
	} else {
		p = price.MultipleSMS
	}
	if p*receiversCount > wallet.Balance {
		return domain.Wallet{}, 0, errors.New("not enough wallet balance")
	}
	return wallet, p * receiversCount, nil
}

func (sh smsHandler) SendSingleSMS(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingSMSRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || request.ReceiverNumber == "" || request.SenderNumber == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if err := ValidateSingleSMSBody(request.SenderNumber, request.ReceiverNumber, request.Content); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, 1)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: request.ReceiverNumber,
		Content:         request.Content,
	}
	err = sh.smsService.SendSMS(smsHistoryRecord)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (sh smsHandler) SendSingleSMSByUsername(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingSMSWithUsernameRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || request.ReceiverUsername == "" || request.SenderNumber == "" || request.PhoneBookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Get the contact
	contact, err := sh.contactService.GetContactByUsername(request.ReceiverUsername)
	if err != nil || contact.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "contact not found"})
	}

	// Check the phone book
	phoneBook, err := sh.phoneBookService.GetPhoneBookById(request.PhoneBookId)
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "this phone book is not for user"})
	}
	if contact.PhoneBookId != phoneBook.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "the contact is not for the given phone book"})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, 1)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: contact.Phone,
		Content:         request.Content,
	}
	err = sh.smsService.SendSMS(smsHistoryRecord)
	if err != nil {

		return c.JSON(http.StatusBadRequest, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (sh smsHandler) SendSinglePeriodSMS(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingPeriodSMSRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || request.ReceiverNumber == "" || request.SenderNumber == "" || request.Period == "" || request.RepeatationCount == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if ValidateSingleSMSBody(request.SenderNumber, request.ReceiverNumber, request.Content) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, request.ReceiverNumber, request.RepeatationCount, sh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS queued"})
}

func (sh smsHandler) SendSinglePeriodSMSByUsername(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingPeriodSMSWithUsernameRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || request.ReceiverUsername == "" || request.SenderNumber == "" || request.Period == "" || request.RepeatationCount == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Get the contact
	contact, err := sh.contactService.GetContactByUsername(request.ReceiverUsername)
	if err != nil || contact.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "contact not found"})
	}

	// Check the phone book
	phoneBook, err := sh.phoneBookService.GetPhoneBookById(request.PhoneBookId)
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "this phone book is not for user"})
	}
	if contact.PhoneBookId != phoneBook.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "the contact is not for the given phone book"})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, contact.Phone, request.RepeatationCount, sh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (sh smsHandler) SendSMSToPhonebooks(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request PhoneBookSMSRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Content == "" || len(request.ReceiverPhoneBooks) == 0 || request.SenderNumber == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	smsHistory, receiversLen, err := sh.smsService.SendSMSToPhonebookIds(domain.SMSHistory{
		Content:      request.Content,
		SenderNumber: request.SenderNumber,
		UserId:       user.ID,
		User:         user,
	}, request.ReceiverPhoneBooks)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, uint(receiversLen))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	err = sh.smsService.SendSMS(smsHistory)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (sh smsHandler) SendPeriodSMSToPhonebooks(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request PhoneBookPeriodSMSRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Content == "" || len(request.ReceiverPhoneBooks) == 0 || request.SenderNumber == "" || request.Period == "" || request.RepeatationCount == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Check for inappropriate words
	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	smsHistory, receiversLen, err := sh.smsService.SendSMSToPhonebookIds(domain.SMSHistory{
		Content:      request.Content,
		SenderNumber: request.SenderNumber,
		UserId:       user.ID,
		User:         user,
	}, request.ReceiverPhoneBooks)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Check the wallet balance and sms price
	wallet, price, err := sh.CheckTheWalletBallence(user, uint(receiversLen)*request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, smsHistory.Content, smsHistory.SenderNumber, smsHistory.ReceiverNumbers, request.RepeatationCount, sh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = sh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}
