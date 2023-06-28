package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PaymentRepository interface {
	Create(input domain.Payment) (domain.Payment, error)
	Update(input domain.Payment) (domain.Payment, error)
	Get(id int) (domain.Payment, error)
}

type PaymentRepositoryImpl struct {
}

func NewPaymentRepository() PaymentRepository {
	return PaymentRepositoryImpl{}
}

func (a PaymentRepositoryImpl) Create(input domain.Payment) (domain.Payment, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (a PaymentRepositoryImpl) Update(input domain.Payment) (domain.Payment, error) {
	var payment domain.Payment
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return payment, err
	}
	_, err = a.Get(int(input.ID))
	if err != nil {
		return payment, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return payment, err
	}

	return payment, nil
}

func (a PaymentRepositoryImpl) Get(id int) (domain.Payment, error) {
	var payment domain.Payment
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&payment, id)

	if err := tx.Error; err != nil {
		return payment, err
	}

	return payment, nil
}
