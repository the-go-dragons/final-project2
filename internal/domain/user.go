package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID              uint      `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	IsLoginRequired bool      `json:"isLoginRequired" gorm:"default:true"`
	IsAdmin         bool      `json:"isAdmin" gorm:"default:false"`
	IsActive        bool      `json:"isActive" gorm:"default:true"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
