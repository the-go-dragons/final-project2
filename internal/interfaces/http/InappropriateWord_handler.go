package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"net/http"
	"strconv"
	"strings"
)

type InappropriateWordHandler struct {
	wordService *usecase.InappropriateWordService
}

func NewInappropriateWordHandler(wordService usecase.InappropriateWordService) InappropriateWordHandler {
	return InappropriateWordHandler{wordService: &wordService}
}

func (n InappropriateWordHandler) Create(c echo.Context) error {
	var req usecase.InappropriateWordDto
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	if len(strings.Trim(req.Word, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid word"})
	}

	dto := usecase.InappropriateWordDto{
		Word: req.Word,
	}

	_, err = n.wordService.Create(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create Inappropriate Word"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n InappropriateWordHandler) Edit(c echo.Context) error {
	var req usecase.InappropriateWordDto
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	if len(strings.Trim(req.Word, " ")) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid word"})
	}

	if req.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid Id"})
	}

	dto := usecase.InappropriateWordDto{
		Word: req.Word,
		ID:   req.ID,
	}

	_, err = n.wordService.Edit(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create Inappropriate Word"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n InappropriateWordHandler) GetAll(c echo.Context) error {
	wordList, err := n.wordService.GetAll()

	if wordList != nil && len(wordList) == 0 {
		return c.JSON(http.StatusOK, wordList)
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't get Inappropriate Words"})
	}

	return c.JSON(http.StatusOK, wordList)
}

func (n InappropriateWordHandler) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "0" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid id"})
	}

	iId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't convert number"})
	}
	err = n.wordService.Delete(uint(iId))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't delete InappropriateWord"})
	}

	return c.JSON(http.StatusOK, Response{Message: "InappropriateWord Deleted Successfully"})
}
