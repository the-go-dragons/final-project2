package usecase

import (
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type PhoneBookService interface {
	CreatePhoneBook(domain.PhoneBook) (domain.PhoneBook, error)
	GetPhoneBookById(uint) (domain.PhoneBook, error)
	GetAllPhoneBooksByUserId(uint) ([]domain.PhoneBook, error)
	GetPhoneBookByUserName(string) ([]domain.PhoneBook, error)
	UpdatePhoneBook(domain.PhoneBook) (domain.PhoneBook, error)
	DeletePhoneBook(uint) error
}

type phoneBookService struct {
	phonebookRepo persistence.PhoneBookRepository
	userRepo      persistence.UserRepository
}

func NewPhoneBook(
	phonebookRepo persistence.PhoneBookRepository,
	userRepo persistence.UserRepository,
) PhoneBookService {
	return phoneBookService{
		phonebookRepo: phonebookRepo,
		userRepo:      userRepo,
	}
}

func (pbs phoneBookService) CreatePhoneBook(input domain.PhoneBook) (domain.PhoneBook, error) {
	// user, err := pbs.userRepo.GetById(input.UserID)
	// if err != nil {
	// 	return domain.PhoneBook{}, err
	// }
	// input.User = user

	return pbs.phonebookRepo.Create(input)
}

func (pbs phoneBookService) GetPhoneBookById(id uint) (domain.PhoneBook, error) {
	return pbs.phonebookRepo.GetById(id)
}

func (pbs phoneBookService) GetAllPhoneBooksByUserId(userId uint) ([]domain.PhoneBook, error) {
	return pbs.phonebookRepo.GetAllByUserId(userId)
}

func (pbs phoneBookService) GetPhoneBookByUserName(username string) ([]domain.PhoneBook, error) {
	user, err := pbs.userRepo.GeByUsername(username)
	if err != nil {
		return make([]domain.PhoneBook, 0), err
	}
	return pbs.phonebookRepo.GetByUser(user)
}

func (pbs phoneBookService) UpdatePhoneBook(input domain.PhoneBook) (domain.PhoneBook, error) {
	// user, err := pbs.userRepo.GetById(input.UserID)
	// if err != nil {
	// 	return domain.PhoneBook{}, err
	// }
	// input.User = user

	return pbs.phonebookRepo.Update(input)
}

func (pbs phoneBookService) DeletePhoneBook(Id uint) error {
	return pbs.phonebookRepo.Delete(Id)
}
