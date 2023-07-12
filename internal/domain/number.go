package domain

import (
	"gorm.io/gorm"
)

type NumberTypeEnum byte

const (
	Sale   NumberTypeEnum = 1
	Rent   NumberTypeEnum = 2
	Public NumberTypeEnum = 3
)

type Number struct {
	gorm.Model
	Phone       string         `json:"phone"`
	Price       uint32         `json:"price"`
	IsAvailable bool           `json:"isAvailable" gorm:"default:true"`
	Type        NumberTypeEnum `json:"type" gorm:"default:1"`
}
