package persistence

import (
	"fmt"
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type NumberRepository interface {
	Create(domain.Number) (domain.Number, error)
	Update(domain.Number) (domain.Number, error)
	Get(uint) (domain.Number, error)
	GetByPhone(string) (domain.Number, error)
	GetDefault() (domain.Number, error)
	GetAllAvailables() ([]domain.Number, error)
	GetUserNumbers(uint) ([]domain.Number, error)
}

type numberRepository struct {
}

func NewNumberRepository() NumberRepository {
	return numberRepository{}
}

func (nr numberRepository) Create(input domain.Number) (domain.Number, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (nr numberRepository) Update(input domain.Number) (domain.Number, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}

func (nr numberRepository) Get(id uint) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&number, id)

	return number, tx.Error
}

func (nr numberRepository) GetByPhone(phone string) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Preload("User").Where("phone = ?", phone).First(&number)

	return number, tx.Error
}

func (nr numberRepository) GetDefault() (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("type = ?", domain.Public).First(&number)

	return number, tx.Error
}

func (nr numberRepository) GetAllAvailables() ([]domain.Number, error) {
	var numbers []domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Table("numbers").
		Joins("FULL JOIN subscriptions ON subscriptions.number_id = numbers.id").
		Where("numbers.user_id IS NULL AND (subscriptions.id IS NULL OR subscriptions.expiration_date < ?)", time.Now()).
		Order("id").
		Find(&numbers)

	return numbers, tx.Error
}

func (nr numberRepository) GetUserNumbers(userId uint) ([]domain.Number, error) {
	var numbers []domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Table("numbers").
		Joins("FULL JOIN subscriptions ON subscriptions.number_id = numbers.id").
		Where("subscriptions.id IS NOT NULL AND subscriptions.expiration_date > ? AND subscriptions.user_id = ?", time.Now(), userId).
		Or("numbers.user_id = ? OR numbers.type = ?", userId, 3).
		Order("id").
		Find(&numbers)
	fmt.Println(tx.Statement.SQL.String())
	return numbers, tx.Error
}
