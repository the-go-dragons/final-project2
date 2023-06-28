package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type TransactionRepository interface {
	Create(input domain.Transaction) (domain.Transaction, error)
	Update(input domain.Transaction) (domain.Transaction, error)
	Get(id int) (domain.Transaction, error)
	GetByWalletID(walletID int) ([]domain.Transaction, error)
}

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() TransactionRepository {
	return TransactionRepositoryImpl{}
}

func (a TransactionRepositoryImpl) Create(input domain.Transaction) (domain.Transaction, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (a TransactionRepositoryImpl) Update(input domain.Transaction) (domain.Transaction, error) {
	var transction domain.Transaction
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return transction, err
	}
	_, err = a.Get(int(input.ID))
	if err != nil {
		return transction, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return transction, err
	}

	return transction, nil
}

func (a TransactionRepositoryImpl) Get(id int) (domain.Transaction, error) {
	var transaction domain.Transaction
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&transaction, id)

	if err := tx.Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (a TransactionRepositoryImpl) GetByWalletID(walletID int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	db, _ := database.GetDatabaseConnection()

	tx := db.Where("wallet_id = ?", walletID).Find(&transactions)

	if err := tx.Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}
