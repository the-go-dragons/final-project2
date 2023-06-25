package domain

import (
	"time"
)

type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at"`
	IsLoginRequired bool      `json:"is_login_required" gorm:"default:true"`
	IsAdmin         bool      `json:"is_admin" gorm:"default:false"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
}
