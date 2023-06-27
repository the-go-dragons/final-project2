package domain

import (
	"time"

	"gorm.io/gorm"
)

type TypeEnum byte

const (
    Sale TypeEnum = 1
    Rent TypeEnum = 2
)

type Number struct {
	gorm.Model
	ID          uint      `json:"id"`
	Phone       int       `json:"phone"`
	IsAvailable bool      `json:"isAvailable" gorm:"default:true"`
	Type        TypeEnum  `json:"type" gorm:"default:1"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
