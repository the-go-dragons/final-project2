package domain

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Balance int64
	UserID  int
}
