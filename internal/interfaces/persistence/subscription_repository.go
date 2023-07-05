package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SubscriptionRepository interface {
	Create(input domain.Subscription) (domain.Subscription, error)
	GetAll() ([]domain.Subscription, error)
	GetByUserId(id uint) (domain.Subscription, error)
	GetByNumber(number domain.Number) (domain.Subscription, error)
}

type SubscriptionRepositoryImpl struct {
}

func NewSubscriptionRepository() SubscriptionRepository {
	return SubscriptionRepositoryImpl{}
}

func (n SubscriptionRepositoryImpl) Create(input domain.Subscription) (domain.Subscription, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (n SubscriptionRepositoryImpl) GetAll() ([]domain.Subscription, error) {
	db, _ := database.GetDatabaseConnection()

	var subscriptions []domain.Subscription

	tx := db.Preload("User").Preload("Number").Find(&subscriptions)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return subscriptions, nil
}

func (a SubscriptionRepositoryImpl) GetByUserId(id uint) (domain.Subscription, error) {
	var wallet domain.Subscription
	db, _ := database.GetDatabaseConnection()

	tx := db.Preload("User").Preload("Number").Where("user_id = ?", id).First(&wallet)

	if err := tx.Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (a SubscriptionRepositoryImpl) GetByNumber(number domain.Number) (domain.Subscription, error) {
	var subscription domain.Subscription
	db, _ := database.GetDatabaseConnection()

	tx := db.Preload("User").Preload("Number").Where("number_id = ?", number.ID).
		// Where("expiration_date > ", time.Now()).
		First(&subscription)

	if err := tx.Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}
