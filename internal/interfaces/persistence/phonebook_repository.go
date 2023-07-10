package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type PhoneBookRepository interface {
	Create(input domain.PhoneBook) (domain.PhoneBook, error)
	Update(input domain.PhoneBook) (domain.PhoneBook, error)
	GetById(id uint) (domain.PhoneBook, error)
	Delete(id uint) error
	GetAll() ([]domain.PhoneBook, error)
	GetByUser(user *domain.User) ([]domain.PhoneBook, error)
}

type PhoneBookRepositoryImpl struct {
}

func NewPhoneBookRepository() PhoneBookRepository {
	return PhoneBookRepositoryImpl{}
}

func (phr PhoneBookRepositoryImpl) Create(input domain.PhoneBook) (domain.PhoneBook, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (phr PhoneBookRepositoryImpl) Update(input domain.PhoneBook) (domain.PhoneBook, error) {
	var phonebook domain.PhoneBook
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return phonebook, err
	}
	_, err = phr.GetById(input.ID)
	if err != nil {
		return phonebook, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return phonebook, err
	}

	return phonebook, nil
}

func (phr PhoneBookRepositoryImpl) GetById(id uint) (domain.PhoneBook, error) {
	var phonebook domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&phonebook, id)

	if err := tx.Error; err != nil {
		return phonebook, err
	}

	return phonebook, nil
}

func (phr PhoneBookRepositoryImpl) Delete(id uint) error {
	var phonebook domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&phonebook, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&phonebook)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (phr PhoneBookRepositoryImpl) GetAll() ([]domain.PhoneBook, error) {
	var phonebooks = make([]domain.PhoneBook, 0)
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&phonebooks)

	tx := db.Debug().Find(&phonebooks)

	if err := tx.Error; err != nil {
		return phonebooks, err
	}

	return phonebooks, nil
}

func (phr PhoneBookRepositoryImpl) GetByUser(user *domain.User) ([]domain.PhoneBook, error) {
	var phonebooks []domain.PhoneBook
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("user_id = ?", user.ID).Find(&phonebooks)

	if err := tx.Error; err != nil {
		return phonebooks, err
	}

	return phonebooks, nil
}
