package usecase

import (
	"errors"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository   *persistence.UserRepository
	walletRepository persistence.WalletRepository
}

func NewUserUsecase(repository *persistence.UserRepository, walletRepository persistence.WalletRepository) *UserUsecase {
	return &UserUsecase{
		userRepository:   repository,
		walletRepository: walletRepository,
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

	createdUser, err := uu.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	wallet := domain.Wallet{
		UserID:  createdUser.ID,
		Balance: 0,
	}

	_, err = uu.walletRepository.Create(wallet)

	return createdUser, err
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

func (uu *UserUsecase) Update(newUser *domain.User) (*domain.User, error) {
	return uu.userRepository.Update(newUser)
}
