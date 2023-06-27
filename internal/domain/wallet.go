package domain

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	User      User      `json:"user"`
	Balance   uint      `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
