package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PriceRepository interface {
	SingltonCreate() (domain.Price, error)
	Update(domain.Price) (domain.Price, error)
}

type priceRepository struct {
}

func NewPriceRepository() PriceRepository {
	return priceRepository{}
}

func (cr priceRepository) SingltonCreate() (domain.Price, error) {
	db, _ := database.GetDatabaseConnection()
	input := domain.Price{}
	tx := db.Where("id = ?", 1).First(&input)
	if input.ID == 0 {
		tx = db.Debug().Create(&domain.Price{})
	}
	return input, tx.Error
}

func (cr priceRepository) Update(input domain.Price) (domain.Price, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}
