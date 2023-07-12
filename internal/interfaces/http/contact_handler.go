package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type ContactHandler interface {
	CreateContact(c echo.Context) error
	GetByPhoneBook(c echo.Context) error
	DeleteContact(c echo.Context) error
}

type contactHandler struct {
	contactService   usecase.ContactService
	phoneBookService usecase.PhoneBookService
}

type ContactData struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

func NewContactHandler(
	contact usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
) ContactHandler {
	return contactHandler{
		contactService:   contact,
		phoneBookService: phoneBookService,
	}
}

func (ch contactHandler) CreateContact(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request ContactData

	// Check the phonebookId from url params
	phonebookId, err := strconv.Atoi(c.Param("phonebookId"))
	if err != nil || phonebookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phonebook id"})
	}

	// Check the request body
	err = c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Phone == "" || request.Username == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.Phone) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	// Check the phone book
	phoneBook, err := ch.phoneBookService.GetPhoneBookById(uint(phonebookId))
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "This phone book is not for user"})
	}

	// Check the dupplication phone in the phone book
	if dupContact, _ := ch.contactService.GetContactByPhone(request.Phone); dupContact.PhoneBookId == phoneBook.ID && dupContact.ID > 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "The phone number already exists in the phonebook"})
	}

	// Check the dupplication username in the phone book
	if dupContact, _ := ch.contactService.GetContactByUsername(request.Username); dupContact.ID > 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "The username already exists"})
	}

	dto := domain.Contact{
		Username:    request.Username,
		Phone:       request.Phone,
		PhoneBookId: uint(phonebookId),
	}

	_, err = ch.contactService.CreateContact(dto)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (ch contactHandler) GetByPhoneBook(c echo.Context) error {
	user := c.Get("user").(domain.User)

	// Check the phonebookId from url params
	phonebookId, err := strconv.Atoi(c.Param("phonebookId"))
	if err != nil || phonebookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phonebook id"})
	}

	// Check the phone book
	phoneBook, err := ch.phoneBookService.GetPhoneBookById(uint(phonebookId))
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "This phone book is not for user"})
	}

	// Get the contacts
	contacts, err := ch.contactService.GetContactByPhoneBookId(phoneBook.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't get the contacts"})
	}

	// Create the response
	response := make([]ContactData, len(contacts))
	for index, contact := range contacts {
		response[index] = ContactData{
			Phone:    contact.Phone,
			Username: contact.Username,
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (ch contactHandler) DeleteContact(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request ContactData

	// Check the phonebookId from url params
	phonebookId, err := strconv.Atoi(c.Param("phonebookId"))
	if err != nil || phonebookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phonebook id"})
	}

	// Check the request body
	err = c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Phone == "" && request.Username == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields, Username or phone is required"})
	}
	if request.Phone != "" && request.Username != "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "Only Username or phone is required"})
	}
	if request.Phone != "" && CheckTheNumberFormat(request.Phone) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	// Check the phone book
	phoneBook, err := ch.phoneBookService.GetPhoneBookById(uint(phonebookId))
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Phone book not found"})
	}
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "This phone book is not for user"})
	}

	// Get the contact
	contact := domain.Contact{}
	if request.Username != "" {
		contact, err = ch.contactService.GetContactByUsername(request.Username)
		if err != nil || contact.ID == 0 {
			return c.JSON(http.StatusBadRequest, Response{Message: "Can't get contact by username"})
		}
	} else {
		contact, err = ch.contactService.GetContactByPhone(request.Phone)
		if err != nil || contact.ID == 0 {
			return c.JSON(http.StatusBadRequest, Response{Message: "Can't get contact by phone"})
		}
	}

	// Delete the contact
	err = ch.contactService.DeleteContact(contact.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Can't delete the contact"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Contact Deleted Successfully"})
}
