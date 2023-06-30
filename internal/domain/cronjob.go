package domain

import (
	"gorm.io/gorm"
)

type CronJob struct {
	gorm.Model
	UserID           uint     `json:"-"`
	User             User     `json:"user"`
	Period           string   `json:"period"`
	RepeatationCount uint     `json:"repeatationCount"`
	Massage          string   `json:"massage"`
	SenderNumber     string   `json:"senderNumber"`
	ReceiverNumbers  []string `json:"receiverNumbers"`
}
