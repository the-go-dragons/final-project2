package domain

import (
	"time"

	"gorm.io/gorm"
)

type NumberTypeEnum byte

const (
    Sale NumberTypeEnum = 1
    Rent NumberTypeEnum = 2
)

type Number struct {
	gorm.Model
	ID          uint      `json:"id"`
	Phone       string    `json:"phone"`
	IsAvailable bool      `json:"isAvailable" gorm:"default:true"`
	Type        NumberTypeEnum  `json:"type" gorm:"default:1"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
