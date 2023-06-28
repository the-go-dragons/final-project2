package domain

import (
	"time"

	"gorm.io/gorm"
)

type SMSTemplate struct {
	gorm.Model
	ID          uint      `json:"id"`
	UserId      uint      `json:"userId"`
	User        User      `json:"user"`
	Text        string    `json:"text"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
} 