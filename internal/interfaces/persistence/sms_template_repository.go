package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SmsTemplateRepository interface {
	Create(domain.SMSTemplate) (domain.SMSTemplate, error)
	GetById(uint) (domain.SMSTemplate, error)
	GetByUserId(uint) ([]domain.SMSTemplate, error)
}

type smsTemplateRepository struct{}

func NewSmsTemplateRepository() SmsTemplateRepository {
	return &smsTemplateRepository{}
}

func (smsR smsTemplateRepository) Create(input domain.SMSTemplate) (domain.SMSTemplate, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Create(&input)

	return input, tx.Error
}

func (smsR smsTemplateRepository) GetById(id uint) (domain.SMSTemplate, error) {
	var smsTemplate domain.SMSTemplate
	db, _ := database.GetDatabaseConnection()

	tx := db.Where("id = ?", id).First(&smsTemplate)

	return smsTemplate, tx.Error
}

func (smsR smsTemplateRepository) GetByUserId(userId uint) ([]domain.SMSTemplate, error) {
	var smsTemplate []domain.SMSTemplate
	db, _ := database.GetDatabaseConnection()
	tx := db.Where("user_id = ?", userId).Find(&smsTemplate)

	return smsTemplate, tx.Error
}
