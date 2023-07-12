package usecase

import (
	"strconv"
	"strings"

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

type SmsPhonebookDto struct {
	PhoneBookdIds []uint      `json:"phoneBookIds"`
	SenderNumber  string      `json:"senderNumber"`
	UserId        uint        `json:"userId"`
	Username      string      `json:"username"`
	User          domain.User `json:"user"`
	Content       string      `json:"content"`
}

func NewSmsService(smsRepo persistence.SmsHistoryRepository,
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

	smsHistory, err := s.CreateSMS(smsHistory)
	if err != nil {
		return err
	}

	return nil
}

func (s SmsServiceImpl) GetSMSHistoryByUserId(userId uint) ([]domain.SMSHistory, error) {
	return s.smsRepo.GetByUserId(userId)
}

func (s SmsServiceImpl) SendToPhonebooks(smsDto SmsPhonebookDto) error {
	contacts, err := s.contactRepo.GetByOfPhoneBookIds(smsDto.PhoneBookdIds)

	if err != nil {
		return err
	}

	for _, phoneBookId := range smsDto.PhoneBookdIds {
		phoneBook, err := s.phonebookRepo.GetById(phoneBookId)

		if err != nil {
			return err
		}

		phoneBookContacts, err := s.contactRepo.GetByPhoneBookId(phoneBook.ID)

		if err != nil {
			return err
		}

		var ids []string
		for _, pb := range phoneBookContacts {
			ids = append(ids, strconv.Itoa(int(pb.ID)))
		}
		phoneBooksContactIds := strings.Join(ids, ",")

		newSmsHistoryRecord := domain.SMSHistory{
			UserId:          smsDto.UserId,
			User:            smsDto.User,
			SenderNumber:    smsDto.SenderNumber,
			ReceiverNumbers: phoneBooksContactIds,
			Content:         smsDto.Content,
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
