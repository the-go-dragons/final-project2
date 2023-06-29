package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type NumberService struct {
	numberRepo  persistence.NumberRepository
}

func NewNumber(
	numberRepo persistence.NumberRepository,
) NumberService {
	return NumberService{
		numberRepo: numberRepo,
	}
}

type NewNumberPayload struct {
	Phone       string                `json:"phone" `
    Type        domain.NumberTypeEnum `json:"type" `
}

func (n NumberService) Create(number NewNumberPayload) (domain.Number, error) {
	now := time.Now()
	numberRecord := domain.Number{
		Phone: number.Phone,
		Type: number.Type,
		IsAvailable: true,
		CreatedAt:  now,
		UpdatedAt: now,
	}

	return n.numberRepo.Create(numberRecord)
}