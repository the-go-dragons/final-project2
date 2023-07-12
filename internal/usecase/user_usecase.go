package usecase

import (
	"errors"
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository   persistence.UserRepository
	walletRepository persistence.WalletRepository
	numberRepository persistence.NumberRepository
	subscriptionRepo persistence.SubscriptionRepository
}

func NewUserUsecase(repository persistence.UserRepository,
	walletRepository persistence.WalletRepository,
	numberRepository persistence.NumberRepository,
	subscriptionRepo persistence.SubscriptionRepository,
) *UserUsecase {
	return &UserUsecase{
		userRepository:   repository,
		walletRepository: walletRepository,
		numberRepository: numberRepository,
		subscriptionRepo: subscriptionRepo,
	}
}

func (uu *UserUsecase) CreateUser(user domain.User) (domain.User, error) {
	// Hash the password
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return user, errors.New("cant hash password")
	}
	user.Password = string(encryptedPassword)

	defaultNumber, err := uu.numberRepository.GetDefault()
	if err == nil {
		user.DefaultNumberID = &defaultNumber.ID
	}
	user.IsActive = true
	user.IsAdmin = false
	user.IsLoginRequired = true

	createdUser, err := uu.userRepository.Create(user)

	if err != nil {
		return user, err
	}

	wallet := domain.Wallet{
		UserID:  createdUser.ID,
		Balance: 0,
	}

	_, err = uu.walletRepository.Create(wallet)

	return createdUser, err
}

func (uu *UserUsecase) GetUserById(id uint) (domain.User, error) {
	return uu.userRepository.GetById(id)
}

func (uu *UserUsecase) GetAll() ([]domain.User, error) {
	return uu.userRepository.GetAll()
}

func (uu *UserUsecase) GetUserByUsername(username string) (domain.User, error) {
	return uu.userRepository.GeByUsername(username)
}

func (uu *UserUsecase) Update(newUser domain.User) (domain.User, error) {
	return uu.userRepository.Update(newUser)
}

func (uu *UserUsecase) UpdateDefaultNumber(userId int, numberId int) (domain.User, error) {
	user, err := uu.userRepository.GetById(uint(userId))
	if err != nil {
		return user, UserNotFound{userId}
	}
	number, err := uu.numberRepository.Get(uint(numberId))
	if err != nil {
		return user, InvalidNumber{int(numberId)}
	}
	if number.Type != domain.Public {
		sub, err := uu.subscriptionRepo.GetByUserId(uint(userId))
		if err != nil || time.Now().After(sub.ExpirationDate) {
			return user, InvalidNumber{int(numberId)}
		}

	}
	user.DefaultNumberID = &number.ID
	return uu.userRepository.Update(user)
}
