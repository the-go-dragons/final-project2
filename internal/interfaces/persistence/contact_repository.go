package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type ContactRepository interface {
	Create(input domain.Contact) (domain.Contact, error)
	GetById(id uint) (domain.Contact, error)
	Delete(id uint) error
	GetByPhoneBookId(phoneBookId uint) ([]domain.Contact, error)
	GetByUsername(username string) (domain.Contact, error)
	GetByPhone(phone string) (domain.Contact, error)
}

type ContactRepositoryImpl struct {
}

func NewContactRepository() ContactRepository {
	return ContactRepositoryImpl{}
}

func (cr ContactRepositoryImpl) Create(input domain.Contact) (domain.Contact, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (cr ContactRepositoryImpl) GetById(id uint) (domain.Contact, error) {
	var Contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&Contact, id)

	if err := tx.Error; err != nil {
		return Contact, err
	}

	return Contact, nil
}

func (cr ContactRepositoryImpl) Delete(id uint) error {
	var Contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&Contact, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&Contact)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (cr ContactRepositoryImpl) GetByPhoneBookId(phoneBookId uint) ([]domain.Contact, error) {
	var contacts []domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone_book_id = ?", phoneBookId).Find(&contacts)

	if err := tx.Error; err != nil {
		return contacts, err
	}

	return contacts, nil
}

func (cr ContactRepositoryImpl) GetByUsername(username string) (domain.Contact, error) {
	var contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("username = ?", username).Find(&contact)

	if err := tx.Error; err != nil {
		return contact, err
	}

	return contact, nil
}

func (cr ContactRepositoryImpl) GetByPhone(phone string) (domain.Contact, error) {
	var contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone = ?", phone).Find(&contact)

	if err := tx.Error; err != nil {
		return contact, err
	}

	return contact, nil
}
