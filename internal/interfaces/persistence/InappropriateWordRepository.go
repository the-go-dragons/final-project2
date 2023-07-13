package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type InappropriateWordRepository interface {
	Create(domain.InappropriateWord) (domain.InappropriateWord, error)
	Update(domain.InappropriateWord) (domain.InappropriateWord, error)
	Get(uint) (domain.InappropriateWord, error)
	Delete(uint) error
	GetAll() ([]domain.InappropriateWord, error)
}

type inappropriateWordRepository struct {
}

func NewInappropriateWordRepository() InappropriateWordRepository {
	return inappropriateWordRepository{}
}

func (iwr inappropriateWordRepository) Create(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (iwr inappropriateWordRepository) Update(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}

func (iwr inappropriateWordRepository) Get(id uint) (domain.InappropriateWord, error) {
	var word domain.InappropriateWord
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&word, id)

	return word, tx.Error
}

func (iwr inappropriateWordRepository) Delete(id uint) error {
	var word domain.InappropriateWord
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&word, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&word)

	return tx.Error
}

func (iwr inappropriateWordRepository) GetAll() ([]domain.InappropriateWord, error) {
	var words []domain.InappropriateWord
	db, _ := database.GetDatabaseConnection()

	tx := db.Find(&words)

	return words, tx.Error
}
