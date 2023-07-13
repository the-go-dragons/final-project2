package domain

import (
	"gorm.io/gorm"
)

type InappropriateWord struct {
	gorm.Model
	Word string `json:"word"`
}
