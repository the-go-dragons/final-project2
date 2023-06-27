package domain

import (
	"time"

	"gorm.io/gorm"
)

type SMSHistory struct {
	gorm.Model
	ID              uint      `json:"id"`
	UserId          uint      `json:"userId"`
	User            User      `json:"user"`
	SenderNumber    string    `json:"senderNumber"`
	ReceiverNumbers string    `json:"receiverNumbers"`
	PhoneBookId     uint      `json:"phoneBookId"`
	PhoneBook       PhoneBook `json:"phoneBook"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
} 