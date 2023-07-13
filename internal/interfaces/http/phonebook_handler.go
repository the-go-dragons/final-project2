package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type PhoneBookHandler struct {
	phonebookService usecase.PhoneBookService
}

func NewPhoneBookHandler(phonebookService usecase.PhoneBookService) PhoneBookHandler {
	return PhoneBookHandler{phonebookService: phonebookService}
}

type NewPhoneBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (pbh PhoneBookHandler) Create(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request NewPhoneBookRequest

	// Check the body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Name == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Missing required fields"})
	}

	// Create the phonebook
	phonebook, err := pbh.phonebookService.CreatePhoneBook(domain.PhoneBook{
		Name:        request.Name,
		Description: request.Description,
		UserID:      user.ID,
	})
	if err != nil || phonebook.ID == 0 {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create phonebook"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (pbh PhoneBookHandler) GetAll(c echo.Context) error {
	user := c.Get("user").(domain.User)
	phonebookList, err := pbh.phonebookService.GetAllPhoneBooksByUserId(user.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't get phone books"})
	}

	return c.JSON(http.StatusOK, phonebookList)
}

func (pbh PhoneBookHandler) Delete(c echo.Context) error {
	user := c.Get("user").(domain.User)
	phonebookId, err := strconv.Atoi(c.Param("id"))
	if err != nil || phonebookId == 0 {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Invalid id"})
	}

	phonebook, err := pbh.phonebookService.GetPhoneBookById(uint(phonebookId))
	if err != nil || phonebook.ID == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Phonebook not found"})
	}

	if phonebook.UserID != user.ID {
		return c.JSON(http.StatusBadRequest, Error{Message: "Phonebook is not for the user"})
	}

	err = pbh.phonebookService.DeletePhoneBook(uint(phonebookId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't delete phone book"})
	}

	return c.JSON(http.StatusOK, Response{Message: "PhoneBook Deleted Successfully"})
}
