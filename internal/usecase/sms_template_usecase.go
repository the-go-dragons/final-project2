package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type SmsTemplateUsecase struct {
	smsTemplateRepository *persistence.SmsTemplateRepository
}

func NewSmsTemplateUsecase(repository *persistence.SmsTemplateRepository) *SmsTemplateUsecase {
	return &SmsTemplateUsecase{
		smsTemplateRepository: repository,
	}
}

func (uu *SmsTemplateUsecase) CreateSMSTemplate(smsTemplate *domain.SMSTemplate) (*domain.SMSTemplate, error) {
	return uu.smsTemplateRepository.Create(smsTemplate)
}

func (uu *SmsTemplateUsecase) GetById(id uint) (*domain.SMSTemplate, error) {
	return uu.smsTemplateRepository.GetById(id)
}

func (uu *SmsTemplateUsecase) GetUserById(id uint) (*domain.SMSTemplate, error) {
	return uu.smsTemplateRepository.GetById(id)
}
