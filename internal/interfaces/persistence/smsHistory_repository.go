package persistence

import (
	"strings"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type SmsHistoryRepository interface {
	Create(domain.SMSHistory) (domain.SMSHistory, error)
	Update(domain.SMSHistory) (domain.SMSHistory, error)
	Get(uint) (domain.SMSHistory, error)
	Delete(uint) error
	GetAll() ([]domain.SMSHistory, error)
	GetByUserId(uint) ([]domain.SMSHistory, error)
	Search(words []string) ([]domain.SMSHistory, error)
}

type smsHistoryRepository struct {
}

func NewSmsHistoryRepository() SmsHistoryRepository {
	return smsHistoryRepository{}
}

func (shr smsHistoryRepository) Create(input domain.SMSHistory) (domain.SMSHistory, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Debug().Create(&input)

	return input, tx.Error
}

func (shr smsHistoryRepository) Update(input domain.SMSHistory) (domain.SMSHistory, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Save(&input)

	return input, tx.Error
}

func (shr smsHistoryRepository) Get(id uint) (domain.SMSHistory, error) {
	var sms domain.SMSHistory
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&sms, id)

	return sms, tx.Error
}

func (shr smsHistoryRepository) Delete(id uint) error {
	var sms domain.SMSHistory
	db, _ := database.GetDatabaseConnection()

	tx := db.First(&sms, id)

	if err := tx.Error; err != nil {
		return err
	}

	tx = tx.Delete(&sms)

	return tx.Error
}

func (shr smsHistoryRepository) GetAll() ([]domain.SMSHistory, error) {
	var smsHistories = make([]domain.SMSHistory, 0)
	db, _ := database.GetDatabaseConnection()
	db = db.Model(&smsHistories)

	tx := db.Debug().Find(&smsHistories)

	return smsHistories, tx.Error
}

func (shr smsHistoryRepository) GetByUserId(userId uint) ([]domain.SMSHistory, error) {
	var sms []domain.SMSHistory
	db, _ := database.GetDatabaseConnection()

	tx := db.Where("user_id = ?", userId).Find(&sms)

	return sms, tx.Error
}
func (shr smsHistoryRepository) Search(words []string) ([]domain.SMSHistory, error) {

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
