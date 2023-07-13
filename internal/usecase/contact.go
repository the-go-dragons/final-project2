package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type ContactService interface {
	CreateContact(domain.Contact) (domain.Contact, error)
	GetContactById(uint) (domain.Contact, error)
	DeleteContact(uint) error
	GetContactByPhoneBookId(uint) ([]domain.Contact, error)
	GetContactByUsername(string) (domain.Contact, error)
	GetContactByPhone(string) (domain.Contact, error)
	GetContactsByListOfPhoneBook([]uint) ([]domain.Contact, error)
}

type contactService struct {
	phonebookRepo    persistence.PhoneBookRepository
	contactRepo      persistence.ContactRepository
	numberRepo       persistence.NumberRepository
	subscriptionRepo persistence.SubscriptionRepository
}

func NewContact(
	phonebookRepo persistence.PhoneBookRepository,
	contactRepo persistence.ContactRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
) ContactService {
	return contactService{
		phonebookRepo:    phonebookRepo,
		contactRepo:      contactRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (cs contactService) CreateContact(contact domain.Contact) (domain.Contact, error) {
	return cs.contactRepo.Create(contact)
}

func (cs contactService) GetContactById(id uint) (domain.Contact, error) {
	return cs.contactRepo.GetById(id)
}

func (cs contactService) DeleteContact(id uint) error {
	return cs.contactRepo.Delete(id)
}

func (cs contactService) GetContactByPhoneBookId(phoneBookId uint) ([]domain.Contact, error) {
	return cs.contactRepo.GetByPhoneBookId(phoneBookId)
}

func (cs contactService) GetContactByUsername(username string) (domain.Contact, error) {
	return cs.contactRepo.GetByUsername(username)
}

func (cs contactService) GetContactByPhone(phone string) (domain.Contact, error) {
	return cs.contactRepo.GetByPhone(phone)
}

func (cs contactService) GetContactsByListOfPhoneBook(phoneBookIds []uint) ([]domain.Contact, error) {
	return cs.contactRepo.GetByOfPhoneBookIds(phoneBookIds)
}
