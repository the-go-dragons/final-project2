package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PriceRepository interface {
	SingltonCreate() (domain.Price, error)
	Update(domain.Price) (domain.Price, error)
}

type priceRepositoryImpl struct {
}

func NewPriceRepository() PriceRepository {
	return priceRepositoryImpl{}
}

func (cr priceRepositoryImpl) SingltonCreate() (domain.Price, error) {
	db, _ := database.GetDatabaseConnection()
	input := domain.Price{}
	db.Where("id = ?", 1).First(&input)
	if input.ID == 0 {
		tx := db.Debug().Create(&input)

		if tx.Error != nil {
			return input, tx.Error
		}
	}
	return input, nil
}

func (cr priceRepositoryImpl) Update(input domain.Price) (domain.Price, error) {
	db, _ := database.GetDatabaseConnection()
	db.Save(&input)
	return input, nil
}
