package usecase

import (
	"errors"
	"strconv"
	"strings"
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

type SmsPhonebookDto struct {
	PhoneBookdIds   []uint      `json:"phoneBookIds"`
	SenderNumber    string      `json:"senderNumber"`
	UserId          uint        `json:"userId"`
	Username        string      `json:"username"`
	User            domain.User `json:"user"`
	Content         string      `json:"content"`
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

	
	newSmsHistoryRecord := domain.SMSHistory{
		UserId:          smsDto.UserId,
		User:            smsDto.User,
		SenderNumber:    smsDto.SenderNumber,
		ReceiverNumbers: smsDto.ReceiverNumbers,
		Content:         smsDto.Content,
	}

	if smsDto.PhoneBookId != 0 {
		phoneBookById, err := s.phonebookRepo.Get(smsDto.PhoneBookId)
		if err != nil {
			return err
		}
		phoneBook = phoneBookById
		newSmsHistoryRecord.PhoneBook = phoneBook
		newSmsHistoryRecord.PhoneBookId = phoneBook.ID
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
		user := &domain.User{
			ID: subscription.UserID,
		}
		phoneBooks, err := s.phonebookRepo.GetByUser(user)
		if err != nil {
			return err
		}
		if len(phoneBooks) == 0 {
			return errors.New("this user has no phonebook")
		}

		newSmsHistoryRecord.PhoneBook = phoneBooks[0]
		newSmsHistoryRecord.PhoneBookId = phoneBooks[0].ID
	}

	now := time.Now()
	newSmsHistoryRecord.CreatedAt = now
	newSmsHistoryRecord.UpdatedAt = now


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

func (s SmsServiceImpl) SendToPhonebooks(smsDto SmsPhonebookDto) error {
	contacts , err := s.contactRepo.GetByPhoneBookIdIn(smsDto.PhoneBookdIds)

	if err != nil {
		return err
	}
	
	for _, phoneBookId := range smsDto.PhoneBookdIds {
		phoneBook, err := s.phonebookRepo.Get(phoneBookId)

		if err != nil {
			return err
		}

		phoneBookContacts, err := s.contactRepo.GetByPhoneBook(&phoneBook)

		if err != nil {
			return err
		}

		var ids []string
		for _, pb := range phoneBookContacts {
			ids = append(ids, strconv.Itoa(int(pb.ID)))
		}
		phoneBooksContactIds := strings.Join(ids, ",")

		now := time.Now()
		newSmsHistoryRecord := domain.SMSHistory{
			UserId:          smsDto.UserId,
			User:            smsDto.User,
			SenderNumber:    smsDto.SenderNumber,
			ReceiverNumbers: phoneBooksContactIds,
			PhoneBook:       phoneBook,
			PhoneBookId:     phoneBookId,
			Content:         smsDto.Content,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		_, err = s.smsRepo.Create(newSmsHistoryRecord)
		if err != nil {
			return err
		}
	}

	for _, contact := range contacts {
		
		// Call the rabbitmq to queue the sms
		smsBody := rabbitmq.SMSBody{
			Sender:    smsDto.SenderNumber,
			Receivers: contact.Phone,
			Massage:   smsDto.Content,
		}

		_ = smsBody

		rabbitmq.NewMassage(smsBody)
	}

	return nil
}