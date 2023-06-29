package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type ContactService struct {
	phonebookRepo persistence.PhoneBookRepository
	contactRepo   persistence.ContactRepository
}

func NewContact(
	phonebookRepo persistence.PhoneBookRepository,
	contactRepo persistence.ContactRepository,
) ContactService {
	return ContactService{
		phonebookRepo: phonebookRepo,
		contactRepo:   contactRepo,
	}
}

type ContactDto struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Phone       string    `json:"phone"`
	PhoneBookId uint      `json:"phoneBookId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (n ContactService) Create(dto ContactDto) (domain.Contact, error) {
	now := time.Now()
	phonebook, err := n.phonebookRepo.Get(dto.PhoneBookId)
	if err != nil {
		return domain.Contact{}, err
	}
	contactRecord := domain.Contact{
		PhoneBookId: dto.PhoneBookId,
		PhoneBook:   phonebook,
		Username:    dto.Username,
		Phone:       dto.Phone,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return n.contactRepo.Create(contactRecord)
}

func (n ContactService) GetById(Id uint) (domain.Contact, error) {
	return n.contactRepo.Get(Id)
}

func (n ContactService) GetAll() ([]domain.Contact, error) {
	return n.contactRepo.GetAll()
}

func (n ContactService) Edit(dto ContactDto) (domain.Contact, error) {
	phonebook, err := n.phonebookRepo.Get(dto.PhoneBookId)
	if err != nil {
		return domain.Contact{}, err
	}
	phonebookRecord := domain.Contact{
		ID:        dto.ID,
		PhoneBook: phonebook,
		Username:  dto.Username,
		Phone:     dto.Phone,
		UpdatedAt: time.Now(),
	}

	return n.contactRepo.Update(phonebookRecord)
}

func (n ContactService) Delete(Id uint) error {
	return n.contactRepo.Delete(Id)
}

func (n ContactService) GetByPhoneBook(phoneBookId uint) ([]domain.Contact, error) {
	phoneBook, err := n.phonebookRepo.Get(phoneBookId)
	if err != nil {
		return make([]domain.Contact, 0), err
	}
	return n.contactRepo.GetByPhoneBook(&phoneBook)
}
