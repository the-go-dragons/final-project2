package domain

import (
	"gorm.io/gorm"
)

type TransactionStatus string

const (
	WITHDRAW TransactionStatus = "Withdraw"
	DEPOSIT  TransactionStatus = "Deposit"
)

type Transaction struct {
	gorm.Model

	UserId         uint              `json:"userId"`
	User           User              `json:"user"`
	WalletID       uint              `json:"walletId"`
	Wallet         int               `json:"wallet"`
	Amount         uint64            `json:"amount"`
	Subscription   Subscription      `json:"subscriptionId"`
	SubscriptionId uint              `json:"subscription"`
	Status         TransactionStatus `json:"status"`
}
