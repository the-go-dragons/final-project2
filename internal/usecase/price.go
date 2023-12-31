package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type PriceService interface {
	Update(domain.Price) (domain.Price, error)
	GetPrice() (domain.Price, error)
}

type priceService struct {
	priceRepository persistence.PriceRepository
}

func NewPriceService(priceRepository persistence.PriceRepository) PriceService {
	return priceService{
		priceRepository: priceRepository,
	}
}

func (ps priceService) Update(input domain.Price) (domain.Price, error) {
	// input2, err := ps.priceRepository.SingltonCreate()
	// if err != nil {
	// 	return input2, err
	// }
	input.ID = 1
	return ps.priceRepository.Update(input)
}

func (ps priceService) GetPrice() (domain.Price, error) {
	return ps.priceRepository.Get()
}
