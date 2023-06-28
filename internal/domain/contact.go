package domain

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID          uint      `json:"id"`
	Username    string `json:"username"`
	Phone       string `json:"phone"`
	PhoneBookId uint `json:"phoneBookId"`
	PhoneBook   PhoneBook `json:"phoneBook"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
