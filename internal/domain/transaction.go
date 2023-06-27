package domain

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	WITHDRAW TransactionStatus = "Withdraw"
	DEPOSIT  TransactionStatus = "Deposit"
)

type Transaction struct {
	gorm.Model
	WalletID        int
	Amount          int64
	TransactionDate time.Time
	Status          TransactionStatus
}
