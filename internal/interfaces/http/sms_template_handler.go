package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type NewSmsTemplateRequest struct {
	Text string `json:"text"`
}

type SmsTemplateHandler struct {
	smsTemplateUseCase *usecase.SmsTemplateUsecase
}

type SmsTemplateResponse struct {
	Message       string `json:"message"`
	SmsTemplateID uint   `json:"smstemplateid"`
}

func NewSmsTemplateHandler(smsTemplateUseCase *usecase.SmsTemplateUsecase) *SmsTemplateHandler {
	return &SmsTemplateHandler{
		smsTemplateUseCase: smsTemplateUseCase,
	}
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

// func (smsh *SmsTemplateHandler) NewSmsWithTemplate(c echo.Context) error {

// }
