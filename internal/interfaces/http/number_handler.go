package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type NumberHandler interface {
	Create(echo.Context) error
	BuyOrRent(echo.Context) error
}

type numberHandler struct {
	numberService usecase.NumberService
	walletService usecase.WalletService
}

func NewNumberHandler(
	numberService usecase.NumberService,
	walletService usecase.WalletService,
) NumberHandler {
	return numberHandler{
		numberService: numberService,
		walletService: walletService,
	}
}

type CreateNubmerRequest struct {
	Phone string                `json:"phone"`
	Price uint32                `json:"price"`
	Type  domain.NumberTypeEnum `json:"type"`
}

type BuyNumberRequest struct {
	NumberId uint `json:"numberId"`
	Months   uint `json:"months"`
}

func (nh numberHandler) Create(c echo.Context) error {
	var request CreateNubmerRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid body request"})
	}
	if request.Phone == "" || request.Price == 0 || request.Type == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.Phone) != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid phone number"})
	}

	// Check number duplicatoin
	number, err := nh.numberService.GetNumberByPhone(request.Phone)
	if err == nil || number.ID != 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Phone number already exists"})
	}

	// Create the number
	number = domain.Number{
		Phone: request.Phone,
		Price: request.Price,
		Type:  request.Type,
	}
	number, err = nh.numberService.CreateNumber(number)
	if err != nil || number.ID == 0 {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (nh numberHandler) BuyOrRent(c echo.Context) error {
	var request BuyNumberRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid body request"})
	}
	if request.NumberId == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Missing required fields"})
	}

	number, err := nh.numberService.GetNumberById(request.NumberId)

	if err != nil || number.ID == 0 {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Number not found"})
	}

	if !number.IsAvailable {
		return c.JSON(http.StatusBadRequest, Response{Message: "Number is not available"})
	}

	if number.Type == 2 && request.Months == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid months"})
	}

	var totalPrice uint32
	var expirationDate time.Time

	if number.Type == 1 {
		totalPrice = number.Price
	} else {
		totalPrice = number.Price * uint32(request.Months)
		expirationDate = time.Now().AddDate(0, int(request.Months), 0)
	}

	user := c.Get("user").(domain.User)

	userWallet, err := nh.walletService.GetByUserId(user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "No wallet found this user"})
	}

	if userWallet.Balance < uint(totalPrice) {
		return c.JSON(http.StatusInternalServerError, Response{Message: "your wallet has not enough balance to pay"})
	}

	_, err = nh.numberService.BuyOrRentNumber(number, user, userWallet, totalPrice, expirationDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't add number: " + err.Error()})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}
