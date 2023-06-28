package usecase

import (
	"fmt"

	"github.com/the-go-dragons/final-project2/internal/domain"
)

type PaymentNotFound struct {
	paymentID int
}

func (o PaymentNotFound) Error() string {
	return fmt.Sprint("Payment not found for payment id:", o.paymentID)
}

type WallertNotFound struct {
	walletId int
}

func (o WallertNotFound) Error() string {
	return fmt.Sprint("Wallet not found for wallet id:", o.walletId)
}

type InvalidBankName struct {
	name string
}

func (i InvalidBankName) Error() string {
	return fmt.Sprint("Invalid bank name: ", i.name)
}

type PaymentNotPaid struct {
	paymentID int
}

func (o PaymentNotPaid) Error() string {
	return fmt.Sprint("Payment is not paid for payment id:", o.paymentID)
}

type PaymentAlreadyApplied struct {
	paymentID int
}

func (o PaymentAlreadyApplied) Error() string {
	return fmt.Sprint("Payment already applied for payment id:", o.paymentID)
}

type InvalidPaymentStatus struct {
	paymentID int
	status    domain.PaymentStatus
}

func (o InvalidPaymentStatus) Error() string {
	return fmt.Sprint("Invalid payment status payment id:", o.paymentID, " with status:", o.status)
}
