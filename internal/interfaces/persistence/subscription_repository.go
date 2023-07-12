package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SubscriptionRepository interface {
	Create(domain.Subscription) (domain.Subscription, error)
	GetAll() ([]domain.Subscription, error)
	GetByUserId(uint) (domain.Subscription, error)
	GetByNumber(domain.Number) (domain.Subscription, error)
}

type subscriptionRepository struct {
}

func NewSubscriptionRepository() SubscriptionRepository {
	return subscriptionRepository{}
}

func (n subscriptionRepository) Create(input domain.Subscription) (domain.Subscription, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (n subscriptionRepository) GetAll() ([]domain.Subscription, error) {
	db, _ := database.GetDatabaseConnection()

	var subscriptions []domain.Subscription

	tx := db.Preload("User").Preload("Number").Find(&subscriptions)

	return subscriptions, tx.Error
}

func (a subscriptionRepository) GetByUserId(id uint) (domain.Subscription, error) {
	var wallet domain.Subscription
	db, _ := database.GetDatabaseConnection()

	tx := db.Preload("User").Preload("Number").Where("user_id = ?", id).First(&wallet)

	return wallet, tx.Error
}

func (a subscriptionRepository) GetByNumber(number domain.Number) (domain.Subscription, error) {
	var subscription domain.Subscription
	db, _ := database.GetDatabaseConnection()

	tx := db.Preload("User").Preload("Number").Where("number_id = ?", number.ID).
		// Where("expiration_date > ", time.Now()).
		First(&subscription)

	return subscription, tx.Error
}
