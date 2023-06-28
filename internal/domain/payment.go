package domain

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PAID    PaymentStatus = "Paid"
	UNPAID  PaymentStatus = "Unpaid"
	APPLIED PaymentStatus = "Applied"
)

type Payment struct {
	gorm.Model

	Amount      uint64        `json:"amount"`
	PaymentDate time.Time     `json:"paymentDate"`
	WalletID    uint          `json:"walletId"`
	Wallet      Wallet        `json:"wallet"`
	Merchant    string        `json:"merchant"`
	Status      PaymentStatus `json:"status"`
}
