package usecase

import (
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type PhoneBookService struct {
	phonebookRepo persistence.PhoneBookRepository
	userRepo      persistence.UserRepository
}

func NewPhoneBook(phonebookRepo persistence.PhoneBookRepository, userRepo *persistence.UserRepository) PhoneBookService {
	return PhoneBookService{
		phonebookRepo: phonebookRepo,
		userRepo:      *userRepo,
	}
}

type PhoneBookDto struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (n PhoneBookService) Create(dto PhoneBookDto) (domain.PhoneBook, error) {
	now := time.Now()
	user, err := n.userRepo.GetById(dto.UserID)
	if err != nil {
		return domain.PhoneBook{}, err
	}
	phonebookRecord := domain.PhoneBook{
		UserID:      dto.UserID,
		User:        *user,
		Name:        dto.Name,
		Description: dto.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return n.phonebookRepo.Create(phonebookRecord)
}

func (n PhoneBookService) GetById(Id uint) (domain.PhoneBook, error) {
	return n.phonebookRepo.Get(Id)
}

func (n PhoneBookService) GetAll() ([]domain.PhoneBook, error) {
	return n.phonebookRepo.GetAll()
}

func (n PhoneBookService) GetByUserName(username string) ([]domain.PhoneBook, error) {
	user, err := n.userRepo.GeByUsername(username)
	if err != nil {
		return make([]domain.PhoneBook, 0), err
	}
	return n.phonebookRepo.GetByUser(user)
}

func (n PhoneBookService) Edit(dto PhoneBookDto) (domain.PhoneBook, error) {
	user, err := n.userRepo.GetById(dto.UserID)
	if err != nil {
		return domain.PhoneBook{}, err
	}
	phonebookRecord := domain.PhoneBook{
		ID:          dto.ID,
		UserID:      dto.UserID,
		Name:        dto.Name,
		User:        *user,
		Description: dto.Description,
		UpdatedAt:   time.Now(),
	}

	return n.phonebookRepo.Update(phonebookRecord)
}

func (n PhoneBookService) Delete(Id uint) error {
	return n.phonebookRepo.Delete(Id)
}
