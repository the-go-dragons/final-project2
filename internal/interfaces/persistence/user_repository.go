package persistence

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	GetById(id uint) (domain.User, error)
	GeByUsername(username string) (domain.User, error)
	Update(user domain.User) (domain.User, error)
	GetAll() ([]domain.User, error)
	Delete(id uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (ur userRepository) Create(user domain.User) (domain.User, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Create(&user)

	return user, tx.Error
}

func (ur userRepository) GetById(id uint) (domain.User, error) {
	var user domain.User
	db, _ := database.GetDatabaseConnection()
	tx := db.Where("id = ?", id).First(&user)

	return user, tx.Error
}

func (ur userRepository) GeByUsername(username string) (domain.User, error) {
	var user domain.User
	db, _ := database.GetDatabaseConnection()
	tx := db.Where("username = ?", username).First(&user)

	return user, tx.Error
}

func (ur userRepository) Update(user domain.User) (domain.User, error) {
	db, _ := database.GetDatabaseConnection()
	tx := db.Save(&user)
	return user, tx.Error
}

func (ur userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	db, _ := database.GetDatabaseConnection()

	tx := db.Debug().Find(&users)

	return users, tx.Error
}

func (ur userRepository) Delete(id uint) error {
	user, err := ur.GetById(id)
	if err != nil {
		return err
	}
	db, _ := database.GetDatabaseConnection()

	tx := db.Delete(&user)

	return tx.Error
}
