package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type InappropriateWordHandler interface {
	CreateInappropriateWord(c echo.Context) error
	GetAll(c echo.Context) error
	Delete(c echo.Context) error
}

type inappropriateWordHandler struct {
	wordService usecase.InappropriateWordService
}

func NewInappropriateWordHandler(
	wordService usecase.InappropriateWordService,
) InappropriateWordHandler {
	return inappropriateWordHandler{
		wordService: wordService,
	}
}

type InappropriateWordResuest struct {
	Word string `json:"word"`
}

func (iwh inappropriateWordHandler) CreateInappropriateWord(c echo.Context) error {
	var request InappropriateWordResuest

	// Check the body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Word == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid body"})
	}

	// Create the inappropriate word record
	word, err := iwh.wordService.Create(domain.InappropriateWord{
		Word: request.Word,
	})
	if err != nil || word.ID == 0 {
		c.JSON(http.StatusBadRequest, Response{Message: "Can't create the word"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (iwh inappropriateWordHandler) GetAll(c echo.Context) error {
	wordList, err := iwh.wordService.GetAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Can't get the word list"})
	}

	return c.JSON(http.StatusOK, wordList)
}

func (iwh inappropriateWordHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid id"})
	}

	err = iwh.wordService.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't delete word"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Word Deleted Successfully"})
}
