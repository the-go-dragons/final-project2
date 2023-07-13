package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type SmsHistoryUsecase struct {
	smsHistoryRepository persistence.SmsHistoryRepository
}

func NewSmsHistoryUsecase(repository persistence.SmsHistoryRepository) SmsHistoryUsecase {
	return SmsHistoryUsecase{
		smsHistoryRepository: repository,
	}
}

func (uu *SmsHistoryUsecase) Search(words []string) ([]domain.SMSHistory, error) {
	return uu.smsHistoryRepository.Search(words)
}
