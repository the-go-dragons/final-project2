package domain

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model

	ID             uint      `json:"id"`
	UserID         uint      `json:"userId"`
	User           User      `json:"user"`
	NumberId       uint      `json:"numberId"`
	Number         Number    `json:"number"`
	Type           TypeEnum  `json:"type" gorm:"default:1"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
} 