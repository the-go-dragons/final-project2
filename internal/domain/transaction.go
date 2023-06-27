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
	
	ID             uint              `json:"id"`
	UserId         uint              `json:"userId"`
	User           User              `json:"user"`
	WalletID       int               `json:"walletId"`
	Wallet         int               `json:"wallet"`
	Amount         uint64            `json:"amount"`
	Subscription   Subscription      `json:"subscriptionId"`
	SubscriptionId uint              `json:"subscription"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Status         TransactionStatus `json:"status"`
}
