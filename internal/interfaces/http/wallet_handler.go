package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type WalletError struct {
	Message string
}

type WalletCharageResponse struct {
	PaymentId uint
}

type WalletFinalizeCharageResponse struct {
	WallertId uint
}
type WalletCharageRequest struct {
	Amount   uint64
	WalletId uint
}

type WalletFinalizeCharageRequest struct {
	PaymentId int
}

type WalletHandler struct {
	wallet *usecase.WalletService
}

func NewWalletHandler(wallet usecase.WalletService) WalletHandler {
	return WalletHandler{wallet: &wallet}
}
func (w WalletHandler) CharageRequest(c echo.Context) error {
	var req WalletCharageRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, WalletError{Message: "Invaild charge request"})
	}
	paymentId, err := w.wallet.ChargeRequest(req.WalletId, req.Amount)
	if err != nil {
		switch err.(type) {
		case usecase.WallertNotFound:
			return c.JSON(http.StatusNotFound, PayResult{Result: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, WalletCharageResponse{paymentId})
}

func (w WalletHandler) FinalizeCharge(c echo.Context) error {
	var walletFChargeReq WalletFinalizeCharageRequest
	err := c.Bind(&walletFChargeReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, WalletError{Message: "Invaild finalize charge request"})
	}

	walletId, err := w.wallet.FinalizeCharge(walletFChargeReq.PaymentId)
	if err != nil {
		switch err.(type) {
		case usecase.PaymentNotFound:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		case *usecase.PaymentNotPaid:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		case usecase.PaymentAlreadyApplied:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		case usecase.InvalidPaymentStatus:
			return c.JSON(http.StatusBadRequest, PayResult{Result: err.Error()})
		default:
			log.Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, WalletFinalizeCharageResponse{walletId})

}
