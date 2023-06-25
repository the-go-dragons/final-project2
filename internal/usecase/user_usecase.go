package usecase

import (
	"errors"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository *persistence.UserRepository
}

func NewUserUsecase(repository *persistence.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: repository,
	}
}

func (uu *UserUsecase) CreateUser(user *domain.User) (*domain.User, error) {
	// Hash the password
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("cant hash password")
	}
	user.Password = string(encryptedPassword)

	user.IsActive = true
	user.IsAdmin = false
	user.IsLoginRequired = true

	return uu.userRepository.Create(user)
}

func (uu *UserUsecase) GetUserById(id uint) (*domain.User, error) {
	return uu.userRepository.GetById(id)
}

func (uu *UserUsecase) GetAll() (*[]domain.User, error) {
	return uu.userRepository.GetAll()
}

func (uu *UserUsecase) GetUserByUsername(username string) (*domain.User, error) {
	return uu.userRepository.GeByUsername(username)
}

func (uu *UserUsecase) UpdateById(id uint, newUser *domain.User) (*domain.User, error) {
	return uu.userRepository.UpdateById(id, newUser)
}
