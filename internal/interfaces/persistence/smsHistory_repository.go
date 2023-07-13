package persistence

import (
	"strings"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SmsHistoryRepository interface {
	Create(input domain.SMSHistory) (domain.SMSHistory, error)
	Update(input domain.SMSHistory) (domain.SMSHistory, error)
	Get(id uint) (domain.SMSHistory, error)
	Delete(id uint) error
	GetAll() ([]domain.SMSHistory, error)
	Search(words []string) ([]domain.SMSHistory, error)
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

func (shr SmsHistoryRepositoryImpl) Search(words []string) ([]domain.SMSHistory, error) {

	var smsHistories = make([]domain.SMSHistory, 0)
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&smsHistories)

	// If no search words are specified, return all SMS history records
	if len(words) == 0 {
		tx := db.Debug().Find(&smsHistories)
		if err := tx.Error; err != nil {
			return smsHistories, err
		}
		return smsHistories, nil
	}

	// Concatenate the input array of words into a single string
	searchString := strings.Join(words, " ")

	// Split the search string into individual words
	searchWords := strings.Fields(searchString)

	// Build the SQL query using the LIKE operator and search words
	query := db.Debug().Where("content LIKE ?", "%"+searchWords[0]+"%")
	for i := 1; i < len(searchWords); i++ {
		query = query.Or("content LIKE ?", "%"+searchWords[i]+"%")
	}

	// Execute the query and retrieve SMS history records
	tx := query.Find(&smsHistories)
	if err := tx.Error; err != nil {
		return smsHistories, err
	}

	return smsHistories, nil
}