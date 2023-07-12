package http

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
)

type SmsHandler struct {
	smsService       *usecase.SmsServiceImpl
	contactService   *usecase.ContactService
	phoneBookService *usecase.PhoneBookService
	wordService      *usecase.InappropriateWordService
}

func NewSmsHandler(
	smsService usecase.SmsServiceImpl,
	contactService usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
	wordService usecase.InappropriateWordService,
) SmsHandler {
	return SmsHandler{
		smsService:       &smsService,
		contactService:   &contactService,
		phoneBookService: &phoneBookService,
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

func (s SmsHandler) SendSingleSMS(c echo.Context) error {
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
	err = s.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (s SmsHandler) SendSingleSMSByUsername(c echo.Context) error {
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
	contact, err := s.contactService.GetContactByUsername(request.ReceiverUsername)
	if err != nil || contact.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "contact not found"})
	}

	// Check the phone book
	phoneBook, err := s.phoneBookService.GetById(request.PhoneBookId)
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "this phone book is not for user"})
	}
	if contact.PhoneBookId != phoneBook.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "the contact is not for the given phone book"})
	}

	err = s.wordService.CheckInappropriateWordsWithRegex(request.Content)
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
	err = s.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (s SmsHandler) SendSinglePeriodSMS(c echo.Context) error {
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
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, request.ReceiverNumber, request.RepeatationCount, s.smsService)

	return c.JSON(http.StatusOK, Response{Message: "SMS queued"})
}

func (s SmsHandler) SendSinglePeriodSMSByUsername(c echo.Context) error {
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
	contact, err := s.contactService.GetContactByUsername(request.ReceiverUsername)
	if err != nil || contact.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "contact not found"})
	}

	// Check the phone book
	phoneBook, err := s.phoneBookService.GetById(request.PhoneBookId)
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
	cronjob.AddNewJob(user, request.Period, request.Content, request.SenderNumber, contact.Phone, request.RepeatationCount, s.smsService)

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (s SmsHandler) SendSMSToPhonebooks(c echo.Context) error {
	var req usecase.SmsPhonebookDto
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	user := c.Get("user").(domain.User)
	if user.ID == 0 {
		return c.JSON(http.StatusNetworkAuthenticationRequired, Response{Message: "Login first"})
	}
	req.User = user
	req.UserId = user.ID

	if !govalidator.Matches(req.SenderNumber, `^(?:\+98)?\d{6,}$`) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid sender number"})
	}

	if len(strings.Trim(req.Content, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid content"})
	}

	if !isUintSlice(req.PhoneBookdIds) || len(req.PhoneBookdIds) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phonebooks"})
	}

	err = s.smsService.SendToPhonebooks(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})

}

func isUintSlice(arr interface{}) bool {
	val := reflect.ValueOf(arr)
	if val.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		if elem.Kind() != reflect.Uint && elem.Kind() != reflect.Uint8 && elem.Kind() != reflect.Uint16 && elem.Kind() != reflect.Uint32 && elem.Kind() != reflect.Uint64 {
			return false
		}
	}
	return true
}
