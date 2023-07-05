package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type SmsHandler struct {
	smsService     *usecase.SmsServiceImpl
	contactService *usecase.ContactService
	wordService    *usecase.InappropriateWordService
}

func NewSmsHandler(smsService usecase.SmsServiceImpl,
	contactService usecase.ContactService,
	wordService *usecase.InappropriateWordService) SmsHandler {
	return SmsHandler{smsService: &smsService, contactService: &contactService, wordService: wordService}
}

func (s SmsHandler) SendSMS(c echo.Context) error {
	var req usecase.SMSHistoryDto
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

	if !govalidator.Matches(req.ReceiverNumbers, `^(?:\+98)?\d{6,}$`) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid receiver number"})
	}

	if strings.EqualFold(req.ReceiverNumbers, req.SenderNumber) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Impossible to send sms to yourself"})
	}

	if len(strings.Trim(req.Content, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid content"})
	}

	err = s.wordService.CheckInappropriateWordsWithRegex(req.Content)
	if err != nil {
		return err
	}

	err = s.smsService.SendSingle(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	err = s.contactService.CreateSmsContact(req.SenderNumber, req.ReceiverNumbers)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create contact"})
	}
	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (s SmsHandler) SendSMSByUsername(c echo.Context) error {
	var req usecase.SMSHistoryDto
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

	if len(strings.Trim(req.Username, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid receiver username"})
	}

	if len(strings.Trim(req.Content, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid content"})
	}

	err = s.wordService.CheckInappropriateWordsWithRegex(req.Content)
	if err != nil {
		return err
	}

	receiverNumber, err := s.smsService.SendSingleByUsername(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms"})
	}

	err = s.contactService.CreateSmsContact(req.SenderNumber, receiverNumber)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create contact"})
	}
	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}
