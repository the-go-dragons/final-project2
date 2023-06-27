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
	Amount      int64
	PaymentDate time.Time
	WalletID    int
	Merchant    string
	Status      PaymentStatus
}
