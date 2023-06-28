package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(input domain.Wallet) (domain.Wallet, error)
	Update(input domain.Wallet) (domain.Wallet, error)
	Get(id int) (domain.Wallet, error)
	ChargeWallet(walletID int, amount int64) error
}

type WalletRepositoryImpl struct {
}

func NewWalletRepository() WalletRepository {
	return WalletRepositoryImpl{}
}

func (a WalletRepositoryImpl) Create(input domain.Wallet) (domain.Wallet, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (a WalletRepositoryImpl) Update(input domain.Wallet) (domain.Wallet, error) {
	var wallet domain.Wallet
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return wallet, err
	}
	_, err = a.Get(int(input.ID))
	if err != nil {
		return wallet, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (a WalletRepositoryImpl) Get(id int) (domain.Wallet, error) {
	var wallet domain.Wallet
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&wallet, id)

	if err := tx.Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (a WalletRepositoryImpl) ChargeWallet(walletID int, amount int64) error {
	var wallet domain.Wallet
	db, _ := database.GetDatabaseConnection()

	tx := db.Model(&wallet).Where("id = ?", walletID).Update("balance", gorm.Expr("balance + ?", amount))

	if err := tx.Error; err != nil {
		return err
	}
	return nil
}
