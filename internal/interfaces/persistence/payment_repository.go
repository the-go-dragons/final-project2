package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PaymentRepository interface {
	Create(domain.Payment) (domain.Payment, error)
	Update(domain.Payment) (domain.Payment, error)
	Get(int) (domain.Payment, error)
}

type paymentRepository struct {
}

func NewPaymentRepository() PaymentRepository {
	return paymentRepository{}
}

func (a paymentRepository) Create(input domain.Payment) (domain.Payment, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (a paymentRepository) Update(input domain.Payment) (domain.Payment, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}

func (a paymentRepository) Get(id int) (domain.Payment, error) {
	var payment domain.Payment
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&payment, id)

	return payment, tx.Error
}
