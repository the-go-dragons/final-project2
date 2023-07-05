package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type NumberService struct {
	numberRepo       persistence.NumberRepository
	walletRepo       persistence.WalletRepository
	subscriptionRepo persistence.SubscriptionRepository
}

func NewNumber(
	numberRepo persistence.NumberRepository,
	walletRepo persistence.WalletRepository,
	subscriptionRepo persistence.SubscriptionRepository,
) NumberService {
	return NumberService{
		numberRepo:       numberRepo,
		walletRepo:       walletRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

type NewNumberPayload struct {
	Phone string                `json:"phone"`
	Price uint32                `json:"price"`
	Type  domain.NumberTypeEnum `json:"type" `
}

func (n NumberService) Create(number NewNumberPayload) (domain.Number, error) {
	numberRecord := domain.Number{
		Phone:       number.Phone,
		Type:        number.Type,
		Price:       number.Price,
		IsAvailable: true,
	}

	return n.numberRepo.Create(numberRecord)
}

func (n NumberService) GetById(Id uint) (domain.Number, error) {
	return n.numberRepo.Get(Id)
}

func (n NumberService) BuyOrRentNumber(number domain.Number, user domain.User, wallet domain.Wallet, totalPrice uint32, expirationDate time.Time) (bool, error) {
	wallet.Balance = wallet.Balance - uint(totalPrice)
	_, err := n.walletRepo.Update(wallet)

	if err != nil {
		return false, err
	}

	subscription := domain.Subscription{
		UserID:         user.ID,
		NumberId:       number.ID,
		Type:           number.Type,
		ExpirationDate: expirationDate,
	}

	_, err = n.subscriptionRepo.Create(subscription)

	if err != nil {
		return false, err
	}

	number.IsAvailable = false
	n.numberRepo.Update(number)

	// n.subscriptionRepo.Create()

	return true, nil
}
