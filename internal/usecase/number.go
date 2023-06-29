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
	Phone       string                `json:"phone"`
	Price       uint32                `json:"price"`
    Type        domain.NumberTypeEnum `json:"type" `
}

func (n NumberService) Create(number NewNumberPayload) (domain.Number, error) {
	now := time.Now()
	numberRecord := domain.Number{
		Phone: number.Phone,
		Type: number.Type,
		Price: number.Price,
		IsAvailable: true,
		CreatedAt:  now,
		UpdatedAt: now,
	}

	return n.numberRepo.Create(numberRecord)
}

func (n NumberService) GetById(Id uint) (domain.Number, error) {
	return n.numberRepo.Get(Id)
}
