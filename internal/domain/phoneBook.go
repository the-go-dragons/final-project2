package domain

import (
	"gorm.io/gorm"
)

type PhoneBook struct {
	gorm.Model
	UserID      uint   `json:"userId"`
	User        User   `json:"user"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
