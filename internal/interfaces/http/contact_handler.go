package http

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/the-go-dragons/final-project2/internal/usecase"
	"net/http"
	"strconv"
)

type ContactHandler struct {
	contactService *usecase.ContactService
}

func NewContactHandler(contact usecase.ContactService) ContactHandler {
	return ContactHandler{contactService: &contact}
}

func (n ContactHandler) Create(c echo.Context) error {
	var req usecase.ContactDto
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

	dto := usecase.ContactDto{
		Username:    req.Username,
		Phone:       req.Phone,
		PhoneBookId: req.PhoneBookId,
	}

	_, err = n.contactService.Create(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

	return c.JSON(http.StatusOK, Response{Message: "Created"})
}

func (n ContactHandler) Edit(c echo.Context) error {
	var req usecase.ContactDto
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

	if req.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "Invalid Id"})
	}

	dto := usecase.ContactDto{
		Username:    req.Username,
		Phone:       req.Phone,
		PhoneBookId: req.PhoneBookId,
		ID:          req.ID,
	}

	_, err = n.contactService.Edit(dto)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Can't create number"})
	}

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
