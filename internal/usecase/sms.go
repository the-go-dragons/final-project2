package usecase

import (
	"errors"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"time"
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
}

type SMSHistoryDto struct {
	ID              uint   `json:"id"`
	UserId          uint   `json:"userId"`
	SenderNumber    string `json:"senderNumber"`
	ReceiverNumbers string `json:"receiverNumbers"`
	PhoneBookId     uint   `json:"phoneBookId"`
	Content         string `json:"content"`
}

func NewSmsService(smsRepo persistence.SmsHistoryRepository,
	userRepo persistence.UserRepository,
	phonebookRepo persistence.PhoneBookRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository) SmsServiceImpl {
	return SmsServiceImpl{
		smsRepo:          smsRepo,
		userRepo:         userRepo,
		phonebookRepo:    phonebookRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s SmsServiceImpl) SendSingle(smsDto SMSHistoryDto) error {
	var user domain.User
	var phoneBook domain.PhoneBook
	var subscription domain.Subscription
	if smsDto.UserId != 0 {
		userById, err := s.userRepo.GetById(smsDto.UserId)
		if err != nil {
			return err
		}
		user = *userById
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
		user = subscription.User
		if user.ID == 0 {
			userById, err := s.userRepo.GetById(subscription.UserID)
			if err != nil {
				return err
			}
			user = *userById
		}
	}
	if smsDto.PhoneBookId != 0 {
		phoneBookById, err := s.phonebookRepo.Get(smsDto.PhoneBookId)
		if err != nil {
			return err
		}
		phoneBook = phoneBookById
	} else {
		if subscription.ID == 0 {
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
		User:            user,
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

	//TODO: call mock api for sending sms if needed.

	return nil
}
