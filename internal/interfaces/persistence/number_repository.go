package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type NumberRepository interface {
	Create(domain.Number) (domain.Number, error)
	Update(domain.Number) (domain.Number, error)
	Get(uint) (domain.Number, error)
	GetByPhone(string) (domain.Number, error)
	GetDefault() (domain.Number, error)
}

type numberRepository struct {
}

func NewNumberRepository() NumberRepository {
	return numberRepository{}
}

func (n numberRepository) Create(input domain.Number) (domain.Number, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (n numberRepository) Update(input domain.Number) (domain.Number, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}

func (a numberRepository) Get(id uint) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&number, id)

	return number, tx.Error
}

func (a numberRepository) GetByPhone(phone string) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone = ?", phone).First(&number)

	return number, tx.Error
}

func (a numberRepository) GetDefault() (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("type = ?", domain.Public).First(&number)

	return number, tx.Error
}
