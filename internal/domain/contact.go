package domain

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Username    string    `json:"username"`
	Phone       string    `json:"phone"`
	PhoneBookId uint      `json:"phoneBookId"`
	PhoneBook   PhoneBook `json:"phoneBook"`
}
