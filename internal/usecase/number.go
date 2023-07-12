package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type NumberService interface {
	CreateNumber(domain.Number) (domain.Number, error)
	GetNumberById(uint) (domain.Number, error)
	BuyOrRentNumber(domain.Number, domain.User, domain.Wallet, uint32, time.Time) (bool, error)
	GetNumberByPhone(string) (domain.Number, error)
}

type numberService struct {
	numberRepository       persistence.NumberRepository
	walletRepository       persistence.WalletRepository
	subscriptionRepository persistence.SubscriptionRepository
}

func NewNumber(
	numberRepository persistence.NumberRepository,
	walletRepository persistence.WalletRepository,
	subscriptionRepository persistence.SubscriptionRepository,
) NumberService {
	return numberService{
		numberRepository:       numberRepository,
		walletRepository:       walletRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (ns numberService) CreateNumber(input domain.Number) (domain.Number, error) {
	return ns.numberRepository.Create(input)
}

func (ns numberService) GetNumberById(id uint) (domain.Number, error) {
	return ns.numberRepository.Get(id)
}

func (ns numberService) BuyOrRentNumber(
	number domain.Number,
	user domain.User,
	wallet domain.Wallet,
	totalPrice uint32,
	expirationDate time.Time,
) (bool, error) {
	wallet.Balance = wallet.Balance - uint(totalPrice)
	_, err := ns.walletRepository.Update(wallet)

	if err != nil {
		return false, err
	}

	subscription := domain.Subscription{
		UserID:         user.ID,
		NumberId:       number.ID,
		Type:           number.Type,
		ExpirationDate: expirationDate,
	}

	_, err = ns.subscriptionRepository.Create(subscription)

	if err != nil {
		return false, err
	}

	number.IsAvailable = false
	ns.numberRepository.Update(number)

	return true, nil
}

func (ns numberService) GetNumberByPhone(phone string) (domain.Number, error) {
	return ns.numberRepository.GetByPhone(phone)
}
