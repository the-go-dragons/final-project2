package domain

import (
	"gorm.io/gorm"
)

type SMSHistory struct {
	gorm.Model
	UserId          uint   `json:"userId"`
	User            User   `json:"user"`
	SenderNumber    string `json:"senderNumber"`
	ReceiverNumbers string `json:"receiverNumbers"`
	Content         string `json:"content"`
}
