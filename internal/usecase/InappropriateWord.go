package usecase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type InappropriateWordService interface {
	Create(input domain.InappropriateWord) (domain.InappropriateWord, error)
	GetById(id uint) (domain.InappropriateWord, error)
	GetAll() ([]domain.InappropriateWord, error)
	Update(input domain.InappropriateWord) (domain.InappropriateWord, error)
	Delete(id uint) error
	CheckInappropriateWordsWithRegex(content string) error
	// CheckInappropriateWords(content string) error
	GetAllInappropriateWords() ([]string, error)
}

type inappropriateWordService struct {
	wordRepository persistence.InappropriateWordRepository
}

func NewInappropriateWord(
	wordRepository persistence.InappropriateWordRepository,
) InappropriateWordService {
	return inappropriateWordService{
		wordRepository: wordRepository,
	}
}

func (iws inappropriateWordService) Create(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	return iws.wordRepository.Create(input)
}

func (iws inappropriateWordService) GetById(id uint) (domain.InappropriateWord, error) {
	return iws.wordRepository.Get(id)
}

func (iws inappropriateWordService) GetAll() ([]domain.InappropriateWord, error) {
	return iws.wordRepository.GetAll()
}

func (iws inappropriateWordService) Update(input domain.InappropriateWord) (domain.InappropriateWord, error) {
	return iws.wordRepository.Update(input)
}

func (iws inappropriateWordService) Delete(id uint) error {
	return iws.wordRepository.Delete(id)
}

func (iws inappropriateWordService) CheckInappropriateWordsWithRegex(content string) error {
	words, err := iws.GetAllInappropriateWords()
	if err != nil {
		return err
	}
	if len(words) == 0 {
		return nil
	}

	matched, err := regexp.MatchString(strings.Join(words, "|"), content)
	if err != nil {
		return err
	}
	if matched {
		return errors.New("inappropriate Word")
	}

	return nil
}

// func (iws inappropriateWordService) CheckInappropriateWords(content string) error {
// 	words, err := iws.GetAllInappropriateWords()
// 	if err != nil {
// 		return err
// 	}

// 	for _, word := range words {
// 		if strings.Contains(content, word) {
// 			return errors.New("inappropriate Word")
// 		}
// 	}

// 	return nil
// }

func (iws inappropriateWordService) GetAllInappropriateWords() ([]string, error) {
	words, err := iws.GetAll()
	if err != nil {
		return make([]string, 0), err
	}

	var result []string

	for _, word := range words {
		result = append(result, word.Word)
	}

	return result, nil
}
