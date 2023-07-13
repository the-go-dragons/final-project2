package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type ContactRepository interface {
	Create(domain.Contact) (domain.Contact, error)
	GetById(uint) (domain.Contact, error)
	Delete(uint) error
	GetByPhoneBookId(uint) ([]domain.Contact, error)
	GetByOfPhoneBookIds([]uint) ([]domain.Contact, error)
	GetByUsername(string) (domain.Contact, error)
	GetByPhone(string) (domain.Contact, error)
}

type contactRepository struct {
}

func NewContactRepository() ContactRepository {
	return contactRepository{}
}

func (cr contactRepository) Create(input domain.Contact) (domain.Contact, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (cr contactRepository) GetById(id uint) (domain.Contact, error) {
	var Contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&Contact, id)

	return Contact, tx.Error
}

func (cr contactRepository) Delete(id uint) error {
	var Contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&Contact, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&Contact)

	return tx.Error
}

func (cr contactRepository) GetByPhoneBookId(phoneBookId uint) ([]domain.Contact, error) {
	var contacts []domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone_book_id = ?", phoneBookId).Find(&contacts)

	return contacts, tx.Error
}

func (cr contactRepository) GetByOfPhoneBookIds(phoneBookIds []uint) ([]domain.Contact, error) {
	var contacts []domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone_book_id in ?", phoneBookIds).Distinct().Find(&contacts)

	return contacts, tx.Error
}

func (cr contactRepository) GetByUsername(username string) (domain.Contact, error) {
	var contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("username = ?", username).First(&contact)

	return contact, tx.Error
}

func (cr contactRepository) GetByPhone(phone string) (domain.Contact, error) {
	var contact domain.Contact
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Where("phone = ?", phone).First(&contact)

	return contact, tx.Error
}
