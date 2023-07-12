package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PhoneBookRepository interface {
	Create(domain.PhoneBook) (domain.PhoneBook, error)
	Update(domain.PhoneBook) (domain.PhoneBook, error)
	GetById(uint) (domain.PhoneBook, error)
	Delete(uint) error
	GetAllByUserId(uint) ([]domain.PhoneBook, error)
	GetByUser(domain.User) ([]domain.PhoneBook, error)
}

type phoneBookRepository struct {
}

func NewPhoneBookRepository() PhoneBookRepository {
	return phoneBookRepository{}
}

func (phr phoneBookRepository) Create(input domain.PhoneBook) (domain.PhoneBook, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (phr phoneBookRepository) Update(input domain.PhoneBook) (domain.PhoneBook, error) {
	db, _ := database.GetDatabaseConnection()

	tx := db.Save(&input)

	return input, tx.Error
}

func (phr phoneBookRepository) GetById(id uint) (domain.PhoneBook, error) {
	var phonebook domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&phonebook, id)

	return phonebook, tx.Error
}

func (phr phoneBookRepository) Delete(id uint) error {
	var phonebook domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&phonebook, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&phonebook)

	return tx.Error
}

func (phr phoneBookRepository) GetAllByUserId(userId uint) ([]domain.PhoneBook, error) {
	var phonebooks []domain.PhoneBook
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&phonebooks)

	tx := db.Where("user_id = ?", userId).Find(&phonebooks)

	return phonebooks, tx.Error
}

func (phr phoneBookRepository) GetByUser(user domain.User) ([]domain.PhoneBook, error) {
	var phonebooks []domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("user_id = ?", user.ID).Find(&phonebooks)

	return phonebooks, tx.Error
}
