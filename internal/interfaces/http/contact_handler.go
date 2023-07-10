package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/internal/usecase"
)

type ContactHandler struct {
	contactService   usecase.ContactService
	phoneBookService usecase.PhoneBookService
}

type ContactRequest struct {
	Username    string `json:"username"`
	Phone       string `json:"phone"`
	PhoneBookId uint   `json:"phonebook_id"`
}

func NewContactHandler(
	contact usecase.ContactService,
	phoneBookService usecase.PhoneBookService,
) ContactHandler {
	return ContactHandler{
		contactService:   contact,
		phoneBookService: phoneBookService,
	}
}

func (n ContactHandler) Create(c echo.Context) error {
	user := c.Get("user").(domain.User)
	var request ContactRequest

	// Check the request body
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}
	if request.Phone == "" || request.Username == "" || request.PhoneBookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Missing required fields"})
	}
	if CheckTheNumberFormat(request.Phone) != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	// Check the phone book
	phoneBook, err := n.phoneBookService.GetById(request.PhoneBookId)
	print(1)
	if err != nil || phoneBook.ID <= 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "phone book not found"})
	}
	print(2)
	if user.ID != phoneBook.UserID {
		return c.JSON(http.StatusBadRequest, Response{Message: "this phone book is not for user"})
	}

	// Check the dupplication phone in the phone book
	if dupContact, _ := n.contactService.GetContactByPhone(request.Phone); dupContact.PhoneBookId == phoneBook.ID && dupContact.ID > 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "the phone number already exists in the phonebook"})
	}

	// Check the dupplication username in the phone book
	if dupContact, _ := n.contactService.GetContactByUsername(request.Username); dupContact.ID > 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "the username already exists"})
	}

	dto := domain.Contact{
		Username:    request.Username,
		Phone:       request.Phone,
		PhoneBookId: request.PhoneBookId,
	}

	_, err = n.contactService.Create(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n ContactHandler) Edit(c echo.Context) error {
	var req ContactRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid data entry"})
	}

	if len(req.Username) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid username"})
	}

	if !govalidator.Matches(req.Phone, `^(?:\+98)?\d{6,}$`) {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid phone number"})
	}

	if req.PhoneBookId == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid PhoneBookId"})
	}

	// dto := ContactRequest{
	// 	Username:    req.Username,
	// 	Phone:       req.Phone,
	// 	PhoneBookId: req.PhoneBookId,
	// }

	// _, err = n.contactService.Edit(dto)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	// }

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n ContactHandler) GetByPhoneBook(c echo.Context) error {
	phonebookId := c.QueryParam("phonebookId")
	if phonebookId == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid phonebookId"})
	}

	iId, err := strconv.Atoi(phonebookId)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}
	contactList, err := n.contactService.GetByPhoneBook(uint(iId))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, contactList)
}

func (n ContactHandler) GetAll(c echo.Context) error {
	contactList, err := n.contactService.GetAll()

	if contactList != nil && len(contactList) == 0 {
		return c.JSON(http.StatusOK, contactList)
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, contactList)
}

func (n ContactHandler) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "0" {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid id"})
	}

	iId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}
	err = n.contactService.Delete(uint(iId))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Contact Deleted Successfully"})
}
