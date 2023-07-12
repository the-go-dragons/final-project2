package usecase

import (
	"errors"
	"strings"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

type SMSService interface {
	CreateSMS(smsHistory domain.SMSHistory) (domain.SMSHistory, error)
	SingleSMS(smsHistory domain.SMSHistory) error
	GetSMSHistoryByUserId(userId uint) ([]domain.SMSHistory, error)
	SendSMSToPhonebookIds(smsHistory domain.SMSHistory, receiverPhoneBookIds []uint) error
}

type smsService struct {
	smsRepo          persistence.SmsHistoryRepository
	userRepo         persistence.UserRepository
	phonebookRepo    persistence.PhoneBookRepository
	numberRepo       persistence.NumberRepository
	subscriptionRepo persistence.SubscriptionRepository
	contactRepo      persistence.ContactRepository
}

func NewSmsService(
	smsRepo persistence.SmsHistoryRepository,
	userRepo persistence.UserRepository,
	phonebookRepo persistence.PhoneBookRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
	contactRepo persistence.ContactRepository,
) SMSService {
	return smsService{
		smsRepo:          smsRepo,
		userRepo:         userRepo,
		phonebookRepo:    phonebookRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
		contactRepo:      contactRepo,
	}
}

func (s smsService) CreateSMS(smsHistory domain.SMSHistory) (domain.SMSHistory, error) {
	return s.smsRepo.Create(smsHistory)
}

func (s smsService) SingleSMS(smsHistory domain.SMSHistory) error {
	// Call the rabbitmq to queue the sms
	smsBody := rabbitmq.SMSBody{
		Sender:    smsHistory.SenderNumber,
		Receivers: smsHistory.ReceiverNumbers,
		Massage:   smsHistory.Content,
	}
	rabbitmq.NewMassage(smsBody)

	// Save the sms history
	smsHistory, err := s.CreateSMS(smsHistory)

	return err
}

func (s smsService) GetSMSHistoryByUserId(userId uint) ([]domain.SMSHistory, error) {
	return s.smsRepo.GetByUserId(userId)
}

func (s smsService) SendSMSToPhonebookIds(smsHistory domain.SMSHistory, receiverPhoneBookIds []uint) error {
	// Get distincted contacts by phone book ids
	contacts, err := s.contactRepo.GetByOfPhoneBookIds(receiverPhoneBookIds)
	if err != nil {
		return err
	}
	if len(contacts) == 0 {
		return errors.New("no contact found")
	}

	// Get the phones
	var receivers []string
	for _, contact := range contacts {
		receivers = append(receivers, contact.Phone)
	}
	smsHistory.ReceiverNumbers = strings.Join(receivers, ",")

	// Call the rabbitmq to queue the sms
	smsBody := rabbitmq.SMSBody{
		Sender:    smsHistory.SenderNumber,
		Receivers: smsHistory.ReceiverNumbers,
		Massage:   smsHistory.Content,
	}
	rabbitmq.NewMassage(smsBody)

	// Save the sms history
	smsHistory, err = s.CreateSMS(smsHistory)

	return err
}
