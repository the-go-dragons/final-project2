package domain

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PAID   PaymentStatus = "Paid"
	UNPAID PaymentStatus = "Unpaid"
)

type Payment struct {
	gorm.Model

	Amount      int64         `json:"amount"`
	PaymentDate time.Time     `json:"paymentDate"`
	WalletID    int           `json:"walletId"`
	Wallet      Wallet        `json:"wallet"`
	Merchant    string	      `json:"merchant"`
	Status      PaymentStatus `json:"status"`
}
