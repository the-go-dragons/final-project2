package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type Error struct {
	Message string
}

type WalletCharageResponse struct {
	PaymentId uint
}

type WalletFinalizeCharageResponse struct {
	WallertId uint
}
type WalletCharageRequest struct {
	Amount uint64
}

type WalletFinalizeCharageRequest struct {
	PaymentId int
}

type WalletHandler struct {
	walletService usecase.WalletService
}

func NewWalletHandler(walletService usecase.WalletService) WalletHandler {
	return WalletHandler{walletService: walletService}
}
func (wh WalletHandler) CharageRequest(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request WalletCharageRequest
	err := c.Bind(&request)
	if err != nil || request.Amount == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invaild amount"})
	}

	wallet, err := wh.walletService.GetByUserId(user.ID)
	if err != nil || wallet.ID == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Can't get the wallet"})
	}

	paymentId, err := wh.walletService.ChargeRequest(wallet.ID, request.Amount)
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

func (wh WalletHandler) FinalizeCharge(c echo.Context) error {
	var walletFChargeReq WalletFinalizeCharageRequest
	err := c.Bind(&walletFChargeReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invaild finalize charge request"})
	}

	walletId, err := wh.walletService.FinalizeCharge(walletFChargeReq.PaymentId)
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
