package domain

import (
	"time"
)

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	CreatedAt       time.Time `json:"createdat"`
	IsLoginRequired bool      `json:"isloginrequired" gorm:"default:false"`
}
