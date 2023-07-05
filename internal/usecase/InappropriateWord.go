package usecase

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/interfaces/persistence"
)

type InappropriateWordService struct {
	wordRepo persistence.InappropriateWordRepository
}

func NewInappropriateWord(
	wordRepo persistence.InappropriateWordRepository,
) InappropriateWordService {
	return InappropriateWordService{
		wordRepo: wordRepo,
	}
}

type InappropriateWordDto struct {
	ID   uint   `json:"id"`
	Word string `json:"word"`
}

func (iws InappropriateWordService) Create(dto InappropriateWordDto) (domain.InappropriateWord, error) {
	now := time.Now()
	wordRecord := domain.InappropriateWord{
		Word:      dto.Word,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return iws.wordRepo.Create(wordRecord)
}

func (iws InappropriateWordService) GetById(Id uint) (domain.InappropriateWord, error) {
	return iws.wordRepo.Get(Id)
}

func (iws InappropriateWordService) GetAll() ([]domain.InappropriateWord, error) {
	return iws.wordRepo.GetAll()
}

func (iws InappropriateWordService) Edit(dto InappropriateWordDto) (domain.InappropriateWord, error) {
	phonebookRecord := domain.InappropriateWord{
		ID:        dto.ID,
		Word:      dto.Word,
		UpdatedAt: time.Now(),
	}

	return iws.wordRepo.Update(phonebookRecord)
}

func (iws InappropriateWordService) Delete(Id uint) error {
	return iws.wordRepo.Delete(Id)
}

func (iws InappropriateWordService) CheckInappropriateWordsWithRegex(content string) error {
	words, err := iws.getAllInappropriateWords()
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

func (iws InappropriateWordService) CheckInappropriateWords(content string) error {
	words, err := iws.getAllInappropriateWords()
	if err != nil {
		return err
	}

	for _, word := range words {
		if strings.Contains(content, word) {
			return errors.New("inappropriate Word")
		}
	}

	return nil
}

func (iws InappropriateWordService) getAllInappropriateWords() ([]string, error) {
	all, err := iws.GetAll()
	if err != nil {
		return make([]string, 0), err
	}

	result := make([]string, len(all))

	for i, word := range all {
		result[i] = word.Word
	}

	return result, nil
}
