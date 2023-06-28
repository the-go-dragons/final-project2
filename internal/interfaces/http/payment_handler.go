package http

import (
	"net/http"

	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type PayError struct {
	Message string
}

type PayResult struct {
	Result string
}

type PaymentHandler struct {
	Payment *usecase.PaymentService
}

func NewPaymentHandler(payment usecase.PaymentService) PaymentHandler {
	return PaymentHandler{Payment: &payment}
}
func (p *PaymentHandler) Pay(c echo.Context) error {
	payemntId, err := strconv.Atoi(c.Param("paymentId"))
	bank := c.QueryParam("bank")
	if bank == "" {
		return c.JSON(http.StatusBadRequest, PayError{Message: "Query parameter 'bank' is required"})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, PayError{Message: "paymentId should be integer"})
	}
	page, err := p.Payment.GetPaymentPage(payemntId, bank)
	if err != nil {
		switch err.(type) {
		case usecase.PaymentNotFound:
			return c.JSON(http.StatusNotFound, PayResult{Result: err.Error()})
		case usecase.InvalidBankName:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, page)
}

func (p *PaymentHandler) Callback(c echo.Context) error {
	form, _ := c.FormParams()
	bank := c.QueryParam("bank")
	if bank == "" {
		return c.JSON(http.StatusBadRequest, PayError{Message: "Query parameter 'bank' is required"})
	}
	result, err := p.Payment.Callback(form, bank)
	if err != nil {
		switch err.(type) {
		case usecase.InvalidBankName:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	if !result.Successful {
		return c.JSON(http.StatusOK, PayResult{Result: "Faild"})
	}
	return c.JSON(http.StatusOK, result)

}
