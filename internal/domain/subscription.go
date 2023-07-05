package domain

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID         uint           `json:"userId"`
	User           User           `json:"user"`
	NumberId       uint           `json:"numberId"`
	Number         Number         `json:"number"`
	Type           NumberTypeEnum `json:"type" gorm:"default:1"`
	ExpirationDate time.Time      `json:"expirationDate"`
}
