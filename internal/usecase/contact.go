package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type ContactService struct {
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
	return ContactService{
		phonebookRepo:    phonebookRepo,
		contactRepo:      contactRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (n ContactService) CreateContact(contact domain.Contact) (domain.Contact, error) {
	return n.contactRepo.Create(contact)
}

func (n ContactService) GetContactById(Id uint) (domain.Contact, error) {
	return n.contactRepo.GetById(Id)
}

func (n ContactService) DeleteContact(Id uint) error {
	return n.contactRepo.Delete(Id)
}

func (n ContactService) GetContactByPhoneBookId(phoneBookId uint) ([]domain.Contact, error) {
	return n.contactRepo.GetByPhoneBookId(phoneBookId)
}

func (n ContactService) GetContactByUsername(username string) (domain.Contact, error) {
	return n.contactRepo.GetByUsername(username)
}

func (n ContactService) GetContactByPhone(phone string) (domain.Contact, error) {
	return n.contactRepo.GetByPhone(phone)
}
