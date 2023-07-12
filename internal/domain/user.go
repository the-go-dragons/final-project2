package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string  `json:"username"`
	Password        string  `json:"password"`
	IsLoginRequired bool    `json:"isLoginRequired" gorm:"default:true"`
	IsAdmin         bool    `json:"isAdmin" gorm:"default:false"`
	IsActive        bool    `json:"isActive" gorm:"default:true"`
	DefaultNumber   *Number `json:"defaultNumber"`
	DefaultNumberID *uint   `json:"defaultNumberID"`
}
