package domain

import (
	"gorm.io/gorm"
)

type SMSTemplate struct {
	gorm.Model
	UserID uint   `json:"-"`
	User   User   `json:"user"`
	Text   string `json:"text"`
}
