package domain

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID  uint `json:"userId"`
	User    User `json:"user"`
	Balance uint `json:"balance"`
}
