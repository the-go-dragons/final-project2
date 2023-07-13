package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
)

type SMSTemplateHandler interface {
	NewSmsTemplate(echo.Context) error
	SmsTemplateList(echo.Context) error
	NewSingleSmsWithTemplate(echo.Context) error
	NewSingleSmsWithUsernameWithTemplate(echo.Context) error
	NewSinglePeriodSmsWithTemplate(echo.Context) error
	NewSinglePeriodSmsWithUsernameWithTemplate(echo.Context) error
	NewPhoneBooksSmsWithTemplate(echo.Context) error
	NewPhoneBooksPeriodSmsWithTemplate(echo.Context) error
}

type smsTemplateHandler struct {
	smsTemplateService usecase.SMSTemplateService
	smsService         usecase.SMSService
	contactService     usecase.ContactService
	phoneBookService   usecase.PhoneBookService
	wordService        usecase.InappropriateWordService
	priceService       usecase.PriceService
}

func NewSmsTemplateHandler(
	smsTemplateService usecase.SMSTemplateService,
	smsService usecase.SMSService,
	contactService usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
	wordService usecase.InappropriateWordService,
	priceService usecase.PriceService,
) SMSTemplateHandler {
	return smsTemplateHandler{
		smsTemplateService: smsTemplateService,
		smsService:         smsService,
		contactService:     contactService,
		phoneBookService:   phoneBookService,
		wordService:        wordService,
		priceService:       priceService,
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

type SinglePeriodSmsWithTemplateRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverNumber   string `json:"receiverNumbers"`
	Content          string `json:"content"`
	TemplateId       uint   `json:"templateId"`
	Period           string `json:"period"`
	RepeatationCount uint   `json:"repeatationCount"`
}

type SinglePeriodSmsWithUsernameWithTemplateRequest struct {
	SenderNumber     string `json:"senderNumber"`
	ReceiverUsername string `json:"receiverUsername"`
	PhoneBookId      uint   `json:"phoneBookId"`
	Content          string `json:"content"`
	TemplateId       uint   `json:"templateId"`
	Period           string `json:"period"`
	RepeatationCount uint   `json:"repeatationCount"`
}

type PhoneBookSmsWithTemplateRequest struct {
	SenderNumber       string `json:"senderNumber"`
	ReceiverPhoneBooks []uint `json:"receiverPhoneBooks"`
	PhoneBookId        uint   `json:"phoneBookId"`
	Content            string `json:"content"`
	TemplateId         uint   `json:"templateId"`
}

type PhoneBookPeriodSmsWithTemplateRequest struct {
	SenderNumber       string `json:"senderNumber"`
	ReceiverPhoneBooks []uint `json:"receiverPhoneBooks"`
	Content            string `json:"content"`
	TemplateId         uint   `json:"templateId"`
	Period             string `json:"period"`
	RepeatationCount   uint   `json:"repeatationCount"`
}

func (smsh smsTemplateHandler) CheckTheWalletBallence(user domain.User, receiversCount uint) (domain.Wallet, uint, error) {
	// Check the wallet balance and sms price
	price, err := smsh.priceService.GetPrice()
	if err != nil || price.ID == 0 {
		return domain.Wallet{}, 0, errors.New("can't get price model")
	}
	wallet, err := smsh.smsService.GetUserWallet(user.ID)
	if err != nil || wallet.ID == 0 {
		return domain.Wallet{}, 0, errors.New("can't get user wallet")
	}
	if price.SingleSMS*receiversCount > wallet.Balance {
		return domain.Wallet{}, 0, errors.New("not enough wallet balance")
	}
	return wallet, price.SingleSMS * receiversCount, nil
}

func (smsh smsTemplateHandler) NewSmsTemplate(c echo.Context) error {
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
	ressmsTemplate, err := smsh.smsTemplateService.CreateSMSTemplate(smsTemplate)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, Response{Message: "Cant create sms template"})
	}

	return c.JSON(http.StatusOK, SmsTemplateResponse{Message: "Created", SmsTemplateID: ressmsTemplate.ID})
}

func (smsh smsTemplateHandler) SmsTemplateList(c echo.Context) error {
	user := c.Get("user").(domain.User)

	templates, err := smsh.smsTemplateService.GetSMSTemplateByUserId(user.ID)
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

func (smsh smsTemplateHandler) NewSingleSmsWithTemplate(c echo.Context) error {
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
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, 1)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: request.ReceiverNumber,
		Content:         content,
	}

	err = smsh.smsService.SendSMS(smsHistoryRecord)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh smsTemplateHandler) NewSingleSmsWithUsernameWithTemplate(c echo.Context) error {
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
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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
	phoneBook, err := smsh.phoneBookService.GetPhoneBookById(request.PhoneBookId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, 1)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Send sms and new sms history
	smsHistoryRecord := domain.SMSHistory{
		UserId:          user.ID,
		User:            user,
		SenderNumber:    request.SenderNumber,
		ReceiverNumbers: contact.Phone,
		Content:         content,
	}

	err = smsh.smsService.SendSMS(smsHistoryRecord)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh smsTemplateHandler) NewSinglePeriodSmsWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SinglePeriodSmsWithTemplateRequest

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

	// Check the template
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, content, request.SenderNumber, request.ReceiverNumber, request.RepeatationCount, smsh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh smsTemplateHandler) NewSinglePeriodSmsWithUsernameWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request SinglePeriodSmsWithUsernameWithTemplateRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || request.ReceiverUsername == "" || request.SenderNumber == "" || request.PhoneBookId == 0 || request.Period == "" || request.RepeatationCount == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Check the template
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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
	phoneBook, err := smsh.phoneBookService.GetPhoneBookById(request.PhoneBookId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, content, request.SenderNumber, contact.Phone, request.RepeatationCount, smsh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh smsTemplateHandler) NewPhoneBooksSmsWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request PhoneBookSmsWithTemplateRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || len(request.ReceiverPhoneBooks) == 0 || request.SenderNumber == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Check the template
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	smsHistory, receiversLen, err := smsh.smsService.SendSMSToPhonebookIds(domain.SMSHistory{
		Content:      content,
		SenderNumber: request.SenderNumber,
		UserId:       user.ID,
		User:         user,
	}, request.ReceiverPhoneBooks)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, uint(receiversLen))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	err = smsh.smsService.SendSMS(smsHistory)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}

func (smsh smsTemplateHandler) NewPhoneBooksPeriodSmsWithTemplate(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request PhoneBookPeriodSmsWithTemplateRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid data entry"})
	}
	if request.Content == "" || len(request.ReceiverPhoneBooks) == 0 || request.SenderNumber == "" || request.Period == "" || request.RepeatationCount == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.SenderNumber) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid sender number"})
	}

	// Check the template
	template, err := smsh.smsTemplateService.GetSMSTemplateById(request.TemplateId)
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

	// Check for inappropriate words
	err = smsh.wordService.CheckInappropriateWordsWithRegex(content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Inappropriate word found"})
	}

	smsHistory, receiversLen, err := smsh.smsService.SendSMSToPhonebookIds(domain.SMSHistory{
		Content:      content,
		SenderNumber: request.SenderNumber,
		UserId:       user.ID,
		User:         user,
	}, request.ReceiverPhoneBooks)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't send sms: " + err.Error()})
	}

	// Check the wallet balance and sms price
	wallet, price, err := smsh.CheckTheWalletBallence(user, uint(receiversLen)*request.RepeatationCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	// Add new cron job
	cronjob.AddNewJob(user, request.Period, smsHistory.Content, smsHistory.SenderNumber, smsHistory.ReceiverNumbers, request.RepeatationCount, smsh.smsService)

	// Change the wallet balance
	wallet.Balance = wallet.Balance - price
	wallet, err = smsh.smsService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't change wallet balance"})
	}

	return c.JSON(http.StatusOK, Response{Message: "SMS Sent"})
}
