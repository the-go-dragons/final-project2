package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type NumberRepository interface {
	Create(input domain.Number) (domain.Number, error)
	Update(input domain.Number) (domain.Number, error)
	Get(id uint) (domain.Number, error)
	GetByPhone(phone string) (domain.Number, error)
	GetDefault() (domain.Number, error)
}

type NumberRepositoryImpl struct {
}

func NewNumberRepository() NumberRepository {
	return NumberRepositoryImpl{}
}

func (n NumberRepositoryImpl) Create(input domain.Number) (domain.Number, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (n NumberRepositoryImpl) Update(input domain.Number) (domain.Number, error) {
	var number domain.Number
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return number, err
	}
	_, err = n.Get(input.ID)
	if err != nil {
		return number, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return number, err
	}

	return number, nil
}

func (a NumberRepositoryImpl) Get(id uint) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&number, id)

	if err := tx.Error; err != nil {
		return number, err
	}

	return number, nil
}

func (a NumberRepositoryImpl) GetByPhone(phone string) (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone = ?", phone).Find(&number)

	if err := tx.Error; err != nil {
		return domain.Number{}, err
	}

	return number, nil
}

func (a NumberRepositoryImpl) GetDefault() (domain.Number, error) {
	var number domain.Number
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("Type = ?", domain.Public).First(&number)

	if err := tx.Error; err != nil {
		return domain.Number{}, err
	}

	return number, nil
}
