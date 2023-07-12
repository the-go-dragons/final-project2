package usecase

import (
	"errors"
	"fmt"
	"time"

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

func (n ContactService) CreateSmsContact(senderNumber string, receiverNumber string) error {
	number, err := n.numberRepo.GetByPhone(senderNumber)
	if err != nil {
		return err
	}
	if number.ID == 0 {
		return errors.New("there is no such a number!")
	}
	subscription, err := n.subscriptionRepo.GetByNumber(number)
	if err != nil {
		return err
	}
	if subscription.ID == 0 || subscription.UserID == 0 {
		return errors.New("this number is assigned to any subscription")
	}

	fmt.Printf("subscription: %v\n", subscription)
	phoneBooks, err := n.phonebookRepo.GetByUser(&subscription.User)
	if err != nil {
		return err
	}
	if len(phoneBooks) == 0 {
		return errors.New("this user has no phonebook")
	}

	phonebookIds := make([]uint, len(phoneBooks))
	for i, phonebook := range phoneBooks {
		phonebookIds[i] = phonebook.ID
	}
	contacts, err := n.contactRepo.GetByPhoneBookIdIn(phonebookIds)
	if err != nil {
		return err
	}
	if len(contacts) == 0 {
		now := time.Now()
		newContact := domain.Contact{
			Phone:       receiverNumber,
			PhoneBook:   phoneBooks[0],
			PhoneBookId: phoneBooks[0].ID,
			CreatedAt:   now,
			UpdatedAt:   now,
			Username:    receiverNumber,
		}

		_, err := n.contactRepo.Create(newContact)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n ContactService) GetContactsByPhonebooks(phoneBookIds []uint) ([]domain.Contact, error) {
	return n.contactRepo.GetByPhoneBookIdIn(phoneBookIds)
}