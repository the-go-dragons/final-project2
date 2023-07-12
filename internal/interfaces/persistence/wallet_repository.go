package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(domain.Wallet) (domain.Wallet, error)
	Update(domain.Wallet) (domain.Wallet, error)
	Get(uint) (domain.Wallet, error)
	ChargeWallet(uint, uint64) error
	GetByUserId(uint) (domain.Wallet, error)
}

type walletRepository struct {
}

func NewWalletRepository() WalletRepository {
	return walletRepository{}
}

func (wr walletRepository) Create(input domain.Wallet) (domain.Wallet, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (wr walletRepository) Update(input domain.Wallet) (domain.Wallet, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Save(&input)

	return input, tx.Error
}

func (wr walletRepository) Get(id uint) (domain.Wallet, error) {
	var wallet domain.Wallet
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&wallet, id)

	return wallet, tx.Error
}

func (wr walletRepository) GetByUserId(id uint) (domain.Wallet, error) {
	var wallet domain.Wallet
	db, _ := database.GetDatabaseConnection()

	tx := db.Where("user_id = ?", id).First(&wallet)

	return wallet, tx.Error
}

func (wr walletRepository) ChargeWallet(walletID uint, amount uint64) error {
	var wallet domain.Wallet
	db, _ := database.GetDatabaseConnection()

	tx := db.Model(&wallet).Where("id = ?", walletID).Update("balance", gorm.Expr("balance + ?", amount))

	return tx.Error
}
