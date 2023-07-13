package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type TransactionRepository interface {
	Create(domain.Transaction) (domain.Transaction, error)
	Update(domain.Transaction) (domain.Transaction, error)
	Get(int) (domain.Transaction, error)
	GetByWalletID(int) ([]domain.Transaction, error)
}

type transactionRepository struct {
}

func NewTransactionRepository() TransactionRepository {
	return transactionRepository{}
}

func (a transactionRepository) Create(input domain.Transaction) (domain.Transaction, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (a transactionRepository) Update(input domain.Transaction) (domain.Transaction, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Save(&input)

	return input, tx.Error
}

func (a transactionRepository) Get(id int) (domain.Transaction, error) {
	var transaction domain.Transaction
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&transaction, id)

	return transaction, tx.Error
}

func (a transactionRepository) GetByWalletID(walletID int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	db, _ := database.GetDatabaseConnection()

	tx := db.Where("wallet_id = ?", walletID).Find(&transactions)

	return transactions, tx.Error
}
