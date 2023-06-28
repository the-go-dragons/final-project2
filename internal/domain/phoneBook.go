package domain

import (
	"time"

	"gorm.io/gorm"
)

type PhoneBook struct {
	gorm.Model
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	User        User      `json:"user"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}