package usecase

import (
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
	UserId          uint        `json:"userId"`
	User            domain.User `json:"user"`
	SenderNumber    string      `json:"senderNumber"`
	ReceiverNumbers string      `json:"receiverNumbers"`
	Content         string      `json:"content"`
}

func NewSmsService(
	smsRepo persistence.SmsHistoryRepository,
	userRepo persistence.UserRepository,
	phonebookRepo persistence.PhoneBookRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
	contactRepo persistence.ContactRepository,
) SmsServiceImpl {
	return SmsServiceImpl{
		smsRepo:          smsRepo,
		userRepo:         userRepo,
		phonebookRepo:    phonebookRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
		contactRepo:      contactRepo,
	}
}

func (s SmsServiceImpl) CreateSMS(smsHistory domain.SMSHistory) (domain.SMSHistory, error) {
	return s.smsRepo.Create(smsHistory)
}

func (s SmsServiceImpl) SingleSMS(smsHistory domain.SMSHistory) error {
	// Call the rabbitmq to queue the sms
	smsBody := rabbitmq.SMSBody{
		Sender:    smsHistory.SenderNumber,
		Receivers: smsHistory.ReceiverNumbers,
		Massage:   smsHistory.Content,
	}
	rabbitmq.NewMassage(smsBody)

	smsHistory, err := s.smsRepo.Create(smsHistory)
	if err != nil {
		return err
	}

	return nil
}
