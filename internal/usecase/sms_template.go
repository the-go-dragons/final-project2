package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type SMSTemplateService interface {
	CreateSMSTemplate(smsTemplate domain.SMSTemplate) (domain.SMSTemplate, error)
	GetSMSTemplateById(id uint) (domain.SMSTemplate, error)
	GetSMSTemplateByUserId(UserId uint) ([]domain.SMSTemplate, error)
}

type smsTemplateService struct {
	smsTemplateRepository persistence.SmsTemplateRepository
}

func NewSmsTemplateService(repository persistence.SmsTemplateRepository) SMSTemplateService {
	return smsTemplateService{
		smsTemplateRepository: repository,
	}
}

func (sts smsTemplateService) CreateSMSTemplate(smsTemplate domain.SMSTemplate) (domain.SMSTemplate, error) {
	return sts.smsTemplateRepository.Create(smsTemplate)
}

func (sts smsTemplateService) GetSMSTemplateById(id uint) (domain.SMSTemplate, error) {
	return sts.smsTemplateRepository.GetById(id)
}

func (sts smsTemplateService) GetSMSTemplateByUserId(UserId uint) ([]domain.SMSTemplate, error) {
	return sts.smsTemplateRepository.GetByUserId(UserId)
}
