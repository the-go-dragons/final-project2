package http

import (
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
}

type smsHandler struct {
	smsService       usecase.SMSService
	contactService   usecase.ContactService
	phoneBookService usecase.PhoneBookService
	wordService      usecase.InappropriateWordService
}

func NewSmsHandler(
	smsService usecase.SMSService,
	contactService usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
	wordService usecase.InappropriateWordService,
) SMSHandler {
	return smsHandler{
		smsService:       smsService,
		contactService:   contactService,
		phoneBookService: phoneBookService,
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
	if ValidateSingleSMSBody(request.SenderNumber, request.ReceiverNumber, request.Content) != nil {
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
	err = sh.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
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

	err = sh.wordService.CheckInappropriateWordsWithRegex(request.Content)
	if err != nil {
		return err
	}

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: contact.Phone,
		Content:         request.Content,
	}
	err = sh.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
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

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, request.ReceiverNumber, request.RepeatationCount, sh.smsService)

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

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, contact.Phone, request.RepeatationCount, sh.smsService)

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

	err = sh.smsService.SendSMSToPhonebookIds(domain.SMSHistory{
		Content:      request.Content,
		SenderNumber: request.SenderNumber,
		UserId:       user.ID,
	}, request.ReceiverPhoneBooks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})

}

// func isUintSlice(arr interface{}) bool {
// 	val := reflect.ValueOf(arr)
// 	if val.Kind() != reflect.Slice {
// 		return false
// 	}
// 	for i := 0; i < val.Len(); i++ {
// 		elem := val.Index(i)
// 		if elem.Kind() != reflect.Uint && elem.Kind() != reflect.Uint8 && elem.Kind() != reflect.Uint16 && elem.Kind() != reflect.Uint32 && elem.Kind() != reflect.Uint64 {
// 			return false
// 		}
// 	}
// 	return true
// }
