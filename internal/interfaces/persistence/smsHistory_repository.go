package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SmsHistoryRepository interface {
	Create(input domain.SMSHistory) (domain.SMSHistory, error)
	Update(input domain.SMSHistory) (domain.SMSHistory, error)
	Get(id uint) (domain.SMSHistory, error)
	Delete(id uint) error
	GetAll() ([]domain.SMSHistory, error)
}

type SmsHistoryRepositoryImpl struct {
}

func NewSmsHistoryRepository() SmsHistoryRepository {
	return SmsHistoryRepositoryImpl{}
}

func (shr SmsHistoryRepositoryImpl) Create(input domain.SMSHistory) (domain.SMSHistory, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	if tx.Error != nil {
		return input, tx.Error
	}

	return input, nil
}

func (shr SmsHistoryRepositoryImpl) Update(input domain.SMSHistory) (domain.SMSHistory, error) {
	var sms domain.SMSHistory
	db, err := database.GetDatabaseConnection()
	if err != nil {
		return sms, err
	}
	_, err = shr.Get(input.ID)
	if err != nil {
		return sms, err
	}
	tx := db.Save(input)
	if err := tx.Error; err != nil {
		return sms, err
	}

	return sms, nil
}

func (shr SmsHistoryRepositoryImpl) Get(id uint) (domain.SMSHistory, error) {
	var sms domain.SMSHistory
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&sms, id)

	if err := tx.Error; err != nil {
		return sms, err
	}

	return sms, nil
}

func (shr SmsHistoryRepositoryImpl) Delete(id uint) error {
	var sms domain.SMSHistory
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&sms, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&sms)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (shr SmsHistoryRepositoryImpl) GetAll() ([]domain.SMSHistory, error) {
	var smsHistories = make([]domain.SMSHistory, 0)
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&smsHistories)

	tx := db.Debug().Find(&smsHistories)

	if err := tx.Error; err != nil {
		return smsHistories, err
	}

	return smsHistories, nil
}
