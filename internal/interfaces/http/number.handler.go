package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type NumberHandler struct {
	number *usecase.NumberService
	wallet *usecase.WalletService
}

type BuyNumberPayload struct {
	NumberId uint `json:"numberId"`
	Months   uint `json:"months"`
}

func NewNumberHandler(number usecase.NumberService, wallet usecase.WalletService) NumberHandler {
	return NumberHandler{number: &number, wallet: &wallet}
}

func (n NumberHandler) Create(c echo.Context) error {
	var req usecase.NewNumberPayload
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid number"})
	}

	if !govalidator.Matches(req.Phone, `^(?:\+98)?\d{6,}$`) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	if !govalidator.IsIn(fmt.Sprintf("%d", req.Type), "1", "2") {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid number type"})
	}

	if !govalidator.IsInt(fmt.Sprintf("%d" , req.Price)) || req.Price ==0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid price"})
	}

	payload := usecase.NewNumberPayload{
		Phone: req.Phone,
		Type: req.Type,
		Price: req.Price,
	}

	_, err = n.number.Create(payload)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}
	
	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n NumberHandler) BuyOrRent(c echo.Context) error {
	var req BuyNumberPayload
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid input"})
	}

	if !govalidator.IsInt(fmt.Sprintf("%d", req.NumberId)) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid number"})
	}

	number, err := n.number.GetById(req.NumberId)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "No number found whit this information"})
	}

	fmt.Printf("number.IsAvailable: %v\n", number.IsAvailable)

	if !number.IsAvailable {
		return c.JSON(http.StatusBadRequest, Response{Message: "Number is not available for buy or rent"})
	}

	if number.Type == 2 && !govalidator.IsInt(fmt.Sprintf("%d", req.Months)) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid months"})
	}

	if number.Type == 2 && req.Months == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid months"})
	}

	var totalPrice uint32
	var expirationDate time.Time

	if number.Type == 1 {
		totalPrice = number.Price
	} else {
		totalPrice = number.Price * uint32(req.Months)
		expirationDate = time.Now().AddDate(0, int(req.Months), 0)
    }

	user := c.Get("user").(domain.User)

	userWallet, err := n.wallet.GetByUserId(user.ID)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "No wallet found this user"})
	}

	if userWallet.Balance < uint(totalPrice) {
		return c.JSON(http.StatusInternalServerError, Response{Message: "your wallet has not enough balance to pay"})
	}

	_, err = n.number.BuyOrRentNumber(number, user, userWallet, totalPrice, expirationDate)

	// _, err = n.number.Create(payload)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create user"})
	}
	
	return c.JSON(http.StatusOK, Response{Message: "Created"})
}