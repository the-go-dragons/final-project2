package usecase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

type SMSService interface {
	CreateSMS(domain.SMSHistory) (domain.SMSHistory, error)
	SendSMS(domain.SMSHistory) error
	GetSMSHistoryByUserId(uint) ([]domain.SMSHistory, error)
	SendSMSToPhonebookIds(domain.SMSHistory, []uint) error
	CheckNumberByUserId(domain.User, string) error
	SendPeriodSMSToPhonebookIds(domain.SMSHistory, []uint) (domain.SMSHistory, error)
	GetUserWallet(uint) (domain.Wallet, error)
	UpdateWallet(domain.Wallet) (domain.Wallet, error)
}

type smsService struct {
	smsRepo          persistence.SmsHistoryRepository
	userRepo         persistence.UserRepository
	phonebookRepo    persistence.PhoneBookRepository
	numberRepo       persistence.NumberRepository
	subscriptionRepo persistence.SubscriptionRepository
	contactRepo      persistence.ContactRepository
	walletRepository persistence.WalletRepository
}

func NewSmsService(
	smsRepo persistence.SmsHistoryRepository,
	userRepo persistence.UserRepository,
	phonebookRepo persistence.PhoneBookRepository,
	numberRepo persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
	contactRepo persistence.ContactRepository,
	walletRepository persistence.WalletRepository,
) SMSService {
	return smsService{
		smsRepo:          smsRepo,
		userRepo:         userRepo,
		phonebookRepo:    phonebookRepo,
		numberRepo:       numberRepo,
		subscriptionRepo: subscriptionRepo,
		contactRepo:      contactRepo,
		walletRepository: walletRepository,
	}
}

func (ss smsService) GetUserWallet(userID uint) (domain.Wallet, error) {
	return ss.walletRepository.GetByUserId(userID)
}

func (ss smsService) UpdateWallet(input domain.Wallet) (domain.Wallet, error) {
	return ss.walletRepository.Update(input)
}

func (s smsService) CreateSMS(smsHistory domain.SMSHistory) (domain.SMSHistory, error) {
	// Replace all digits with length more than 4 with asterisks
	re := regexp.MustCompile("\\d{5,}")
	smsHistory.Content = re.ReplaceAllStringFunc(smsHistory.Content, func(match string) string {
		asterisks := ""
		for i := 0; i < len(match); i++ {
			asterisks += "*"
		}
		return asterisks
	})

	// Create the SMS history record in the repository
	return s.smsRepo.Create(smsHistory)
}

func (ss smsService) CheckNumberByUserId(user domain.User, phone string) error {
	number, err := ss.numberRepo.GetByPhone(phone)
	if err != nil || number.ID == 0 {
		return errors.New("number not found")
	}

	if number.Type == domain.Public {
		return nil
	} else if number.Type == domain.Sale {
		// If the number is for sale and is not for the user, return false
		if number.User == nil || *number.UserID != user.ID {
			return errors.New("number is not for the user")
		}
	} else if number.Type == domain.Rent {
		subscriptions, _ := ss.subscriptionRepo.GetNotExpiredByNumber(number.ID)
		if len(subscriptions) == 0 {
			return errors.New("number is not for the user")
		}
		for _, sub := range subscriptions {
			if sub.UserID == user.ID {
				return nil
			}
		}
		return errors.New("number is not available")
	} else {
		return errors.New("invalid number type")
	}
	return nil
}

func (s smsService) SendSMS(smsHistory domain.SMSHistory) error {
	// Check the sender number
	err := s.CheckNumberByUserId(smsHistory.User, smsHistory.SenderNumber)
	if err != nil {
		return err
	}

	// Check the wallet balance

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

func (s smsService) GetSMSHistoryByUserId(userId uint) ([]domain.SMSHistory, error) {
	return s.smsRepo.GetByUserId(userId)
}

func (s smsService) SendSMSToPhonebookIds(smsHistory domain.SMSHistory, receiverPhoneBookIds []uint) error {
	// Check the sender number
	err := s.CheckNumberByUserId(smsHistory.User, smsHistory.SenderNumber)
	if err != nil {
		return err
	}

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

	s.SendSMS(smsHistory)

	return err
}

func (s smsService) SendPeriodSMSToPhonebookIds(smsHistory domain.SMSHistory, receiverPhoneBookIds []uint) (domain.SMSHistory, error) {
	// Check the sender number
	err := s.CheckNumberByUserId(smsHistory.User, smsHistory.SenderNumber)
	if err != nil {
		return smsHistory, err
	}

	// Get distincted contacts by phone book ids
	contacts, err := s.contactRepo.GetByOfPhoneBookIds(receiverPhoneBookIds)
	if err != nil {
		return smsHistory, err
	}
	if len(contacts) == 0 {
		return smsHistory, errors.New("no contact found")
	}

	// Get the phones
	var receivers []string
	for _, contact := range contacts {
		receivers = append(receivers, contact.Phone)
	}
	smsHistory.ReceiverNumbers = strings.Join(receivers, ",")

	return smsHistory, err
}
