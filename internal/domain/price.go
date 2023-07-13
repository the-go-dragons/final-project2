package domain

import (
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	SingleSMS   uint `json:"singleSms"`
	MultipleSMS uint `json:"multipleSms"`
}
