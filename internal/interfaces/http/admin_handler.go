package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type AdminHandler interface {
	DisableUser(echo.Context) error
	ChangePricing(echo.Context) error
	GetSMSHistoryByUserId(echo.Context) error
	UsersList(echo.Context) error
}

type adminHandler struct {
	userUsecase  usecase.UserUsecase
	priceUsecase usecase.PriceService
	smsService   usecase.SMSService
}

func NewAdminHandler(
	userUsecase usecase.UserUsecase,
	priceUsecase usecase.PriceService,
	smsService usecase.SMSService,
) AdminHandler {
	return adminHandler{
		userUsecase:  userUsecase,
		priceUsecase: priceUsecase,
		smsService:   smsService,
	}
}

type ChangePricingRequest struct {
	SingleSMS   uint `json:"single_sms"`
	MultipleSMS uint `json:"multiple_sms"`
}

type ChangePricingResponse struct {
	domain.Price
}

type SMSHistoryResponse struct {
	Count        int                 `json:"count"`
	SMSHistories []domain.SMSHistory `json:"sms_histories"`
}

func (ah adminHandler) UsersList(c echo.Context) error {
	_ = c.Get("user").(domain.User)

	users, err := ah.userUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't get users list"})
	}

	return c.JSON(http.StatusOK, users)
}

func (ah adminHandler) DisableUser(c echo.Context) error {
	admin := c.Get("user").(domain.User)

	// Check the user id from url params
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid user id"})
	}
	user, err := ah.userUsecase.GetUserById(uint(userId))
	if err != nil || user.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "User not found"})
	}
	if !user.IsActive {
		return c.JSON(http.StatusBadRequest, Response{Message: "User already disabled"})
	}
	if user.ID == admin.ID {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can not disable your self"})
	}
	if user.IsAdmin {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can not disable other admins"})
	}

	// Update the user
	user.IsActive = false
	ah.userUsecase.Update(user)

	return c.JSON(http.StatusOK, Response{Message: "User disabled successfully"})
}

func (ah adminHandler) ChangePricing(c echo.Context) error {
	_ = c.Get("user").(domain.User)
	var request ChangePricingRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.SingleSMS == 0 || request.MultipleSMS == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}

	// Update the price
	price := domain.Price{
		SingleSMS:   request.SingleSMS,
		MultipleSMS: request.MultipleSMS,
	}
	price, err = ah.priceUsecase.Update(price)
	if err != nil || price.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can not update the price: " + err.Error()})
	}

	return c.JSON(http.StatusOK, ChangePricingResponse{price})
}

func (ah adminHandler) GetSMSHistoryByUserId(c echo.Context) error {
	_ = c.Get("user").(domain.User)

	// Check the user id from url params
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid user id"})
	}
	user, err := ah.userUsecase.GetUserById(uint(userId))
	if err != nil || user.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "User not found"})
	}

	// Get sms histories
	smsHistories, err := ah.smsService.GetSMSHistoryByUserId(user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Could not get the sms histories: " + err.Error()})
	}

	return c.JSON(http.StatusOK, SMSHistoryResponse{Count: len(smsHistories), SMSHistories: smsHistories})
}
