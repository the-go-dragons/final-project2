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

func (n ContactService) Create(contact domain.Contact) (domain.Contact, error) {
	return n.contactRepo.Create(contact)
}

func (n ContactService) GetById(Id uint) (domain.Contact, error) {
	return n.contactRepo.Get(Id)
}

func (n ContactService) GetAll() ([]domain.Contact, error) {
	return n.contactRepo.GetAll()
}

func (n ContactService) Edit(contact domain.Contact) (domain.Contact, error) {
	return n.contactRepo.Update(contact)
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

func (n ContactService) GetContactByUsername(username string) (domain.Contact, error) {
	return n.contactRepo.GetByUsername(username)
}

func (n ContactService) GetContactByPhone(phone string) (domain.Contact, error) {
	return n.contactRepo.GetByPhone(phone)
}
