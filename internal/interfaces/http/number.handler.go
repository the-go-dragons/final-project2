package http

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type NumberHandler struct {
	number *usecase.NumberService
}

func NewNumberHandler(number usecase.NumberService) NumberHandler {
	return NumberHandler{number: &number}
}

func (n NumberHandler) Create(c echo.Context) error {
	var req usecase.NewNumberPayload
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, WalletError{Message: "Invalid number"})
	}

	if !govalidator.Matches(req.Phone, `^(?:\+98|0)?9\d{9}$`) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	if !govalidator.IsIn(fmt.Sprintf("%d", req.Type), "1", "2") {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid number type"})
	}

	payload := usecase.NewNumberPayload{
		Phone: req.Phone,
		Type: req.Type,
	}

	_, err = n.number.Create(payload)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create user"})
	}
	
	return c.JSON(http.StatusOK, Response{Message: "Created"})
}