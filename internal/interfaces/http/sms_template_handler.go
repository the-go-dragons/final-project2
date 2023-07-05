package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type SmsTemplateHandler struct {
	smsTemplateUseCase *usecase.SmsTemplateUsecase
	smsService         *usecase.SmsServiceImpl
	contactService     *usecase.ContactService
	phoneBookService   *usecase.PhoneBookService
}

func NewSmsTemplateHandler(
	smsTemplateUseCase *usecase.SmsTemplateUsecase,
	smsService usecase.SmsServiceImpl,
	contactService usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
) *SmsTemplateHandler {
	return &SmsTemplateHandler{
		smsTemplateUseCase: smsTemplateUseCase,
		smsService:         &smsService,
		contactService:     &contactService,
		phoneBookService:   &phoneBookService,
	}
}

type NewSmsTemplateRequest struct {
	Text string `json:"text"`
}

type SmsTemplateResponse struct {
	Message       string `json:"message"`
	SmsTemplateID uint   `json:"smstemplateid"`
}

type SingleSmsWithTemplateRequest struct {
	SenderNumber   string `json:"senderNumber"`
	ReceiverNumber string `json:"receiverNumbers"`
	Content        string `json:"content"`
	TemplateId     uint   `json:"templateId"`
}

type SingleSmsWithUsernameWithTemplateRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverUsername string `json:"receiverUsername"`
	PhoneBookId      uint   `json:"phoneBookId"`
	Content          string `json:"content"`
	TemplateId       uint   `json:"templateId"`
}

func (smsh *SmsTemplateHandler) NewSmsTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request NewSmsTemplateRequest

	// Check the body data
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid body request"})

	}
	if request.Text == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}

	// Check the count of inputs
	count := strings.Count(request.Text, "%s")
	if count <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Must have at least one argument with %s"})
	}

	// Create the sms template
	smsTemplate := domain.SMSTemplate{
		UserID: user.ID,
		Text:   request.Text,
	}
	ressmsTemplate, err := smsh.smsTemplateUseCase.CreateSMSTemplate(&smsTemplate)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Cant create sms template"})
	}

	return c.JSON(http.StatusOK, SmsTemplateResponse{Message: "Created", SmsTemplateID: ressmsTemplate.ID})
}

func (smsh *SmsTemplateHandler) SmsTemplateList(c echo.Context) error {
	user := c.Get("user").(domain.User)

	templates, err := smsh.smsTemplateUseCase.GetByUserId(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Cant get sms template list"})
	}

	var response []SmsTemplateResponse

	for _, template := range templates {
		response = append(response, SmsTemplateResponse{
			Message:       template.Text,
			SmsTemplateID: template.ID,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (smsh *SmsTemplateHandler) NewSingleSmsWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingleSmsWithTemplateRequest

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

	// Check the template
	template, err := smsh.smsTemplateUseCase.GetById(request.TemplateId)
	if err != nil || template.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Template not found"})
	}
	if template.UserID != user.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "The selected template is not for the user"})
	}

	// Make the content with the template
	slices := strings.Split(string(request.Content), "%")
	interfaceSlice := make([]interface{}, len(slices))

	for i, v := range slices {
		interfaceSlice[i] = v
	}
	content := fmt.Sprintf(template.Text, interfaceSlice...)

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: request.ReceiverNumber,
		Content:         content,
	}

	err = smsh.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh *SmsTemplateHandler) NewSingleSmsWithUsernameWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SingleSmsWithUsernameWithTemplateRequest

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

	// Check the template
	template, err := smsh.smsTemplateUseCase.GetById(request.TemplateId)
	if err != nil || template.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Template not found"})
	}
	if template.UserID != user.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "The selected template is not for the user"})
	}

	// Get the contact
	contact, err := smsh.contactService.GetContactByUsername(request.ReceiverUsername)
	if err != nil || contact.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "contact not found"})
	}

	// Check the phone book
	phoneBook, err := smsh.phoneBookService.GetById(request.PhoneBookId)
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "this phone book is not for user"})
	}
	if contact.PhoneBookId != phoneBook.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "the contact is not for the given phone book"})
	}

	// Make the content with the template
	slices := strings.Split(string(request.Content), "%")
	interfaceSlice := make([]interface{}, len(slices))

	for i, v := range slices {
		interfaceSlice[i] = v
	}
	content := fmt.Sprintf(template.Text, interfaceSlice...)

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: contact.Phone,
		Content:         content,
	}

	err = smsh.smsService.SingleSMS(smsHistoryRecord)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}
