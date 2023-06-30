package usecase

import (
	"errors"
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

type SmsService interface {
	SendSingle(smsDto SMSHistoryDto) error
}

type SmsServiceImpl struct {
	smsRepo          persistence.SmsHistoryRepository
	userRepo         persistence.UserRepository
	phonebookRepo    persistence.PhoneBookRepository
	numberRepo       persistence.NumberRepository
	subscriptionRepo persistence.SubscriptionRepository
	contactRepo      persistence.ContactRepository
}

type SMSHistoryDto struct {
	ID              uint        `json:"id"`
	UserId          uint        `json:"userId"`
	SenderNumber    string      `json:"senderNumber"`
	ReceiverNumbers string      `json:"receiverNumbers"`
	PhoneBookId     uint        `json:"phoneBookId"`
	Content         string      `json:"content"`
	Username        string      `json:"username"`
	User            domain.User `json:"user"`
}

func NewSmsService(smsRepo persistence.SmsHistoryRepository,
	userRepo persistence.UserRepository,
	phonebookRepo persistence.PhoneBookRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
	contactRepo persistence.ContactRepository) SmsServiceImpl {
	return SmsServiceImpl{
		smsRepo:          smsRepo,
		userRepo:         userRepo,
		phonebookRepo:    phonebookRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
		contactRepo:      contactRepo,
	}
}

func (s SmsServiceImpl) SendSingle(smsDto SMSHistoryDto) error {
	var phoneBook domain.PhoneBook
	if smsDto.PhoneBookId != 0 {
		phoneBookById, err := s.phonebookRepo.Get(smsDto.PhoneBookId)
		if err != nil {
			return err
		}
		phoneBook = phoneBookById
	} else {
		number, err := s.numberRepo.GetByPhone(smsDto.SenderNumber)
		if err != nil {
			return err
		}
		if number.ID == 0 {
			return errors.New("there is no such a number")
		}
		subscription, err := s.subscriptionRepo.GetByNumber(number)
		if err != nil {
			return err
		}
		if subscription.ID == 0 || subscription.UserID == 0 {
			return errors.New("this number is assigned to any subscription")
		}
		phoneBooks, err := s.phonebookRepo.GetByUser(&subscription.User)
		if err != nil {
			return err
		}
		if len(phoneBooks) == 0 {
			return errors.New("this user has no phonebook")
		}
	}

	now := time.Now()

	newSmsHistoryRecord := domain.SMSHistory{
		UserId:          smsDto.UserId,
		User:            smsDto.User,
		SenderNumber:    smsDto.SenderNumber,
		ReceiverNumbers: smsDto.ReceiverNumbers,
		PhoneBook:       phoneBook,
		PhoneBookId:     smsDto.PhoneBookId,
		Content:         smsDto.Content,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err := s.smsRepo.Create(newSmsHistoryRecord)
	if err != nil {
		return err
	}

	// Call the rabbitmq to queue the sms
	smsBody := rabbitmq.SMSBody{
		Sender:    newSmsHistoryRecord.SenderNumber,
		Receivers: newSmsHistoryRecord.ReceiverNumbers,
		Massage:   newSmsHistoryRecord.Content,
	}
	rabbitmq.NewMassage(smsBody)

	return nil
}

func (s SmsServiceImpl) SendSingleByUsername(smsDto SMSHistoryDto) (string, error) {
	var newSmsHistoryRecord domain.SMSHistory
	receiverNumber := ""
	contact, err := s.contactRepo.GetByUsername(smsDto.Username)
	if err != nil {
		return receiverNumber, err
	}
	phoneBooks, err := s.phonebookRepo.GetByUser(&smsDto.User)
	if err != nil {
		return receiverNumber, err
	}
	if len(phoneBooks) == 0 {
		return receiverNumber, errors.New("this user has no phonebook")
	}

	if contact.ID != 0 {
		receiverNumber = contact.Phone
	} else {
		receiverUser, err := s.userRepo.GeByUsername(smsDto.Username)
		if err != nil {
			return receiverNumber, err
		}
		subscription, err := s.subscriptionRepo.GetByUserId(receiverUser.ID)
		if err != nil {
			return receiverNumber, err
		}
		if subscription.ID == 0 || subscription.Number.ID == 0 {
			return receiverNumber, errors.New("no number found for this user")
		}

		receiverNumber = subscription.Number.Phone
	}

	now := time.Now()
	newSmsHistoryRecord = domain.SMSHistory{
		UserId:          smsDto.UserId,
		User:            smsDto.User,
		SenderNumber:    smsDto.SenderNumber,
		ReceiverNumbers: receiverNumber,
		PhoneBook:       phoneBooks[0],
		PhoneBookId:     phoneBooks[0].ID,
		Content:         smsDto.Content,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = s.smsRepo.Create(newSmsHistoryRecord)
	if err != nil {
		return receiverNumber, err
	}

	// Call the rabbitmq to queue the sms
	smsBody := rabbitmq.SMSBody{
		Sender:    newSmsHistoryRecord.SenderNumber,
		Receivers: newSmsHistoryRecord.ReceiverNumbers,
		Massage:   newSmsHistoryRecord.Content,
	}
	rabbitmq.NewMassage(smsBody)

	return receiverNumber, nil
}
