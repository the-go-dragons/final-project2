package persistence

import (
	"errors"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SmsTemplateRepository struct{}

func NewSmsTemplateRepository() *SmsTemplateRepository {
	return &SmsTemplateRepository{}
}

func (smsR *SmsTemplateRepository) Create(smsTemplate *domain.SMSTemplate) (*domain.SMSTemplate, error) {
	db, _ := database.GetDatabaseConnection()
	result := db.Create(&smsTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	return smsTemplate, nil
}

func (smsR *SmsTemplateRepository) GetById(id uint) (*domain.SMSTemplate, error) {
	smsTemplate := new(domain.SMSTemplate)
	db, _ := database.GetDatabaseConnection()
	db.Where("id = ?", id).First(&smsTemplate)
	if smsTemplate.ID == 0 {
		return nil, errors.New("SMSTemplate not found")
	}
	return smsTemplate, nil
}

func (smsR *SmsTemplateRepository) GetByUserId(userId uint) ([]domain.SMSTemplate, error) {
	var smsTemplate []domain.SMSTemplate
	db, _ := database.GetDatabaseConnection()
	db.Where("user_id = ?", userId).Find(&smsTemplate)
	return smsTemplate, nil
}
