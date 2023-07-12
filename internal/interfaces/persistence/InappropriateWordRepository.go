package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type InappropriateWordRepository interface {
	Create(input domain.InappropriateWord) (domain.InappropriateWord, error)
	Update(input domain.InappropriateWord) (domain.InappropriateWord, error)
	Get(id uint) (domain.InappropriateWord, error)
	Delete(id uint) error
	GetAll() ([]domain.InappropriateWord, error)
}

type InappropriateWordRepositoryImpl struct {
}

func NewInappropriateWordRepository() InappropriateWordRepository {
	return InappropriateWordRepositoryImpl{}
}

func (iwr InappropriateWordRepositoryImpl) Create(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (iwr InappropriateWordRepositoryImpl) Update(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	var word domain.InappropriateWord
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return word, err
	}
	_, err = iwr.Get(input.ID)
	if err != nil {
		return word, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return word, err
	}

	return word, nil
}

func (iwr InappropriateWordRepositoryImpl) Get(id uint) (domain.InappropriateWord, error) {
	var word domain.InappropriateWord
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&word, id)

	if err := tx.Error; err != nil {
		return word, err
	}

	return word, nil
}

func (iwr InappropriateWordRepositoryImpl) Delete(id uint) error {
	var word domain.InappropriateWord
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&word, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&word)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (iwr InappropriateWordRepositoryImpl) GetAll() ([]domain.InappropriateWord, error) {
	var words = make([]domain.InappropriateWord, 0)
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&words)

	tx := db.Debug().Find(&words)

	if err := tx.Error; err != nil {
		return words, err
	}

	return words, nil
}
