package domain

import (
	"gorm.io/gorm"
	"time"
)

type InappropriateWord struct {
	gorm.Model
	ID        uint      `json:"id"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
